package odinX

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/offer365/example/etcd/dao"
	"github.com/offer365/example/etcd/embedder"
	"github.com/offer365/example/qrcode"
	"github.com/offer365/odin/asset"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/utils"

	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

const (
	restfulUser               = "admin"
	defaultKey                = "default"
	statusOk            int   = 200
	statusMethodErr     int32 = 405
	statusChkLicenseErr int32 = iota + 440
	statusPutLicenseErr
	statusClearLicenseErr
	statusUntiedAppErr
	statusGetLicenseErr
)

var (
	Self *Node
	// auth          *Authentication
	store         dao.Store
	device        embedder.Embed
	confWhiteList = map[string]string{"default": "", "members": ""} // 在白名单的配置无法删除
	_assetPath    string
	secrets       = gin.H{"admin": nil}
	Lock          sync.Mutex
	gs            *grpc.Server
)

// 释放静态资源
func RestoreAsset() {
	// 解压 静态文件的位置
	if runtime.GOOS == "linux" {
		_assetPath = "/usr/share/.asset/.temp/"
	} else {
		_assetPath = "./"
	}
	// go get -u github.com/jteeuwen/go-bindata/...
	// 重新生成静态资源在项目的根目录下 go-bindata -o=asset/asset.go -pkg=asset html/... static/...
	dirs := []string{"html", "static"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets(_assetPath, dir); err != nil {
			log.Sugar.Error("restore assets failed. error: ", err)
			_ = os.RemoveAll(filepath.Join(_assetPath, dir))
			continue
		}
	}
}

func Main() {
	var (
		err   error
		ready = make(chan struct{})
	)
	// 生产模式
	gin.SetMode(gin.ReleaseMode)
	// 添加basic auth 认证
	secrets[Cfg.GRpcUser] = Cfg.GRpcPwd
	Self = NewNode(Cfg.NodeName, Cfg.NodeAddr)
	RestoreAsset()
	device = embedder.NewEmbed()
	if err = device.Init(Cfg.EmbedCtx,
		embedder.WithName(Cfg.EmbedName),
		embedder.WithDir(Cfg.EmbedDir),
		embedder.WithClientAddr(Cfg.EmbedClientAddr),
		embedder.WithPeerAddr(Cfg.EmbedPeerAddr),
		embedder.WithClusterToken(Cfg.EmbedClusterToken),
		embedder.WithClusterState(Cfg.EmbedClusterState),
		embedder.WithCluster(Cfg.EmbedCluster),
		embedder.WithLogger(log.Sugar),
	); err != nil {
		log.Sugar.Error("init embed server failed. error: ", err)
	}

	go func() { // 运行etcd
		if err = device.Run(ready); err != nil {
			log.Sugar.Error("run embed server error. ", err)
			return
		}
	}()
	select {
	case <-ready: // 待etcd Ready 运行其他服务
		err = device.SetAuth(Cfg.EmbedAuthUser, Cfg.EmbedAuthPwd)
		if err != nil {
			log.Sugar.Error("set auth embed server failed. error: ", err)
		}
		close(ready)
		Server()
	}
}

func Server() {
	var (
		err error
	)
	// 客户端连接
	store = dao.NewStore()
	if err = store.Init(Cfg.EtcdCliCtx,
		dao.WithAddr(Cfg.EtcdCliAddr),
		dao.WithUsername(Cfg.EmbedAuthUser),
		dao.WithPassword(Cfg.EmbedAuthPwd),
		dao.WithTimeout(Cfg.EtcdCliTimeout),
	); err != nil {
		log.Sugar.Error("init store failed. error: ", err)
	}
	// 从etcd加载license
	if err := loadLic(); err != nil {
		log.Sugar.Error("init license failed. error: ", err)
	}

	// 间隔1分钟更新授权
	go func() {
		ticker := time.Tick(1 * time.Minute) // 1分钟
		// expr := cronexpr.MustParse("* * * * *")
		for range ticker {
			// now := time.Date()
			// next := expr.Next(now)
			// time.AfterFunc(next.Sub(now), func() {
			// time.AfterFunc(time.Second, func() {})
			// 如果是主就更新授权
			if device.IsLeader() {
				log.Sugar.Infof("%s is Leader. ip:%s", Self.Attrs.Name, Self.Attrs.Addr)
				if err := ResetLicense(); err != nil {
					log.Sugar.Error("reset license failed. error: ", err)
				}
			}
		}
	}()
	// 监听授权变化
	go WatchLicense()
	go RunAPI(Cfg.GRpcListen)
	go WebServer()
	AllNodeGRpcClient(Cfg.GRpcAllNode)
	DefaultConf()
	signalChan := make(chan os.Signal)
	done := make(chan struct{}, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, os.Kill)
	// 资源回收
	go func() {
		<-signalChan
		ClientConns.Range(func(key, value interface{}) bool {
			cli, ok := value.(*grpc.ClientConn)
			if ok {
				cli.Close()
			}
			return true
		})
		gs.Stop()
		device.Close()
		store.Close()
		close(signalChan)
		done <- struct{}{}
	}()
	// 阻塞主进程
	<-done
	// <-make(chan struct{})
	// <- (chan int)(nil)
}

// Load authorization when launching the program
func loadLic() (err error) {
	var (
		byt []byte
		lic *License
	)
	if byt, err = GetLicense(); err != nil {
		log.Sugar.Error("get license failed. error: ", err)
	}

	if byt == nil || len(byt) == 0 {
		lic = new(License)
		StoreLic(lic)
		return
	}
	// 如果此位置发生错误，则硬盘中的许可证可能无效（数据损坏）
	if lic, err = Str2lic(string(byt)); err != nil {
		log.Sugar.Error("load license error. may be data corruption: ", err)
		if err = DelLicense(); err != nil {
			log.Sugar.Error("load license error. when delete license. ", err)
		}
		lic = new(License)
		StoreLic(lic)
		return
	}

	if lic != nil {
		now := time.Now().Unix()
		// 检查服务器当前时间是否小于license 时间
		if now < lic.Update {
			log.Sugar.Errorf("check license update time error. %d %d", now, lic.Update)
			if err = DelLicense(); err != nil {
				log.Sugar.Error("check license error. when delete license. ", err)
			}
			lic = new(License)
			StoreLic(lic)
			return
		}
		// 检查license 的生存周期
		if (now-lic.Generate)/60 < lic.LifeCycle {
			log.Sugar.Errorf("check license life cycle error. %d %d", (now-lic.Generate)/60, lic.LifeCycle)
			if err = DelLicense(); err != nil {
				log.Sugar.Error("check license error. when delete license. ", err)
			}
			lic = new(License)
			StoreLic(lic)
			return
		}
		// 检查硬件信息是否匹配
		if lic.Devices[Self.Attrs.Name] != Self.Attrs.Hwmd5 {
			log.Sugar.Error("check license hard ware md5 error.")
			if err = DelLicense(); err != nil {
				log.Sugar.Error("check license error. when delete license. ", err)
			}
			lic = new(License)
			StoreLic(lic)
			return
		}

	} else {
		lic = new(License)
	}
	StoreLic(lic)
	return
}

func RunAPI(addr string) {
	var err error
	gs, err = NodeGRpcServer()
	if err != nil {
		log.Sugar.Error(err)
		return
	}
	RegisterStaterServer(gs, Self)
	RegisterAuthorizeServer(gs, Author)
	ws := apiServer()
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gs.ServeHTTP(w, r) // grpc server
		} else {
			ws.ServeHTTP(w, r) // gin api server
		}
		return
	})
	listener, err := NewTlsListen([]byte(Cfg.GRpcServerCrt), []byte(Cfg.GRpcServerKey), []byte(Cfg.GRpcCaCrt), addr)
	if err != nil {
		log.Sugar.Error(err)
		return
	}
	err = http.Serve(listener, handle)
	if err != nil {
		log.Sugar.Error(err)
		return
	}
}

func NewTlsListen(crt, key, ca []byte, addr string) (net.Listener, error) {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Sugar.Error(err)
		return nil, err
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = errors.New("failed to append ca certs")
		log.Sugar.Error(err)
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.RequireAndVerifyClientCert, // NOTE: 这是可选的!
		ClientCAs:          certPool,
		InsecureSkipVerify: false,
		Rand:               rand.Reader,
		Time:               time.Now,
		NextProtos:         []string{"http/1.1", http2.NextProtoTLS},
	}
	return tls.Listen("tcp", addr, tlsConfig)
}

// gin 路由
func WebServer() {
	r := gin.New()
	r.Use(gin.Recovery()) // Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	// store := cookie.NewStore([]byte("secret"))
	// store.Options(sessions.Options{MaxAge:0})
	// r.Use(sessions.Sessions("odin",store))
	r.LoadHTMLGlob(_assetPath + "html/*")
	r.Static("/static", _assetPath+"static")
	r.StaticFile("/favicon.ico", _assetPath+"static/favicon.ico")
	// api 路由组 basicauth 认证
	api := r.Group("/odin/api/v1", gin.BasicAuth(gin.Accounts{
		restfulUser: Cfg.WebPwd,
	}))

	// 序列号
	api.Any("/server/code", SerialNumAPI)
	// 序列号二维码
	api.GET("/server/qr-code", QrCodeAPI)
	// 授权码
	api.Any("/server/license", LicenseAPI)
	// 注销二维码
	api.GET("/server/qr-license", QrLicenseAPI)
	// 节点状态
	api.GET("/server/nodes", NodeStatusAPI)
	// 解绑接口
	api.POST("/server/untied/:app/:id", UntiedAppApi)
	// 配置接口 kv 存储
	api.Any("/client/conf/*name", ConfAPI)
	// 客户端在线接口
	api.GET("/client/online/*app", CliOnlineAPI)
	// web 交互
	api.POST("/web/login", LoginAPI)
	// 帮助
	api.GET("/server/help", HelpAPI)

	r.Any("", func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		r.HandleContext(c)
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "首页",
		})
	})
	log.Sugar.Info(r.Run(Cfg.WebListen))
}

// gin 路由
func apiServer() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery()) // Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	// store := cookie.NewStore([]byte("secret"))
	// store.Options(sessions.Options{MaxAge:0})
	// r.Use(sessions.Sessions("odin",store))
	api := r.Group("/odin/api/v1", gin.BasicAuth(gin.Accounts{
		Cfg.GRpcUser: Cfg.GRpcPwd,
	}))
	// 客户端接口
	api.Any("/client/auth", ClientAPI)
	return r
}

// 客户端Api
func ClientAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if pwd, ok := secrets[user]; ok && pwd == Cfg.GRpcPwd {
		// app := c.Param("app")
		// id := c.Param("id")
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Sugar.Error("ready request body error: ", err)
			return
		}
		req := new(Request)
		err = json.Unmarshal(data, req)
		if err != nil {
			return
		}
		switch c.Request.Method {
		case "POST": // 认证
			resp, err := Author.Auth(context.TODO(), req)
			if err != nil {
				log.Sugar.Error(resp, err)
				return
			}
			c.JSON(statusOk, resp)
			return
		case "PUT": // 心跳
			resp, err := Author.KeepLine(context.TODO(), req)
			if err != nil {
				log.Sugar.Error(resp, err)
				return
			}
			c.JSON(statusOk, resp)
			return
		case "DELETE": // 关闭
			resp, err := Author.OffLine(context.TODO(), req)
			if err != nil {
				log.Sugar.Error(resp, err)
				return
			}
			c.JSON(statusOk, resp)
			return
		default:
			c.JSON(statusOk, &Response{Code: statusMethodErr, Data: nil, Msg: "method error",})
		}
	}
}

type online struct {
	ID   string `json:"id"`
	Info string `json:"info"`
}

type stat struct {
	ID     string `json:"id"`
	Addr   string `json:"addr"`
	Online string `json:"online"`
}

type conf struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// 序列号Api
func SerialNumAPI(c *gin.Context) {
	var (
		code string
		err  error
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	code = "auth error"
	if _, ok := secrets[user]; ok {
		switch c.Request.Method {
		case "GET":
			code, err = GetSerialNum()
		case "POST":
			code, err = ResetSerialNum()
		default:
			code = "method error"
		}
	}
	if err != nil {
		code = err.Error()
		log.Sugar.Error(code)
	}
	c.JSON(statusOk, gin.H{
		"code": statusOk,
		"data": map[string]string{"key": code, "date": time.Now().Format("2006-01-02 15:04:05")},
		"msg":  "在授权成功前请勿重启进程或系统，否则序列号将变更。请保证机器硬件和系统时间正确(误差5分钟内)，否则可能会导致进程异常或者授权失效。",
	})
}

// 序列号二维码Api
func QrCodeAPI(c *gin.Context) {
	var (
		code string
		err  error
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	code = "auth error"
	if _, ok := secrets[user]; ok {
		code, err = GetSerialNum()
	}
	if err != nil {
		code = err.Error()
		log.Sugar.Error(code)
	}

	// c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))  //对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Type", "image/jpeg") // 不能去掉
	buf := bytes.NewBuffer(make([]byte, 0))
	err = qrcode.NewQrEncode([]byte(code), buf)
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="qr-code.jpg"`,
	}
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "image/jpeg", buf, extraHeaders)

	// 使用文件实现
	// c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))
	// c.Writer.Header().Add("Content-Type", "application/octet-stream")
	// c.File(AssetPath+"qr-code.jpg")
}

// 注销二维码Api
func QrLicenseAPI(c *gin.Context) {
	var (
		code string
		err  error
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	code = "auth error"
	if _, ok := secrets[user]; ok {
		code, err = GetClearLicense()
	}
	if err != nil {
		code = err.Error()
		log.Sugar.Error(code)
	}

	// c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))  //对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Type", "image/jpeg") // 不能去掉
	buf := bytes.NewBuffer(make([]byte, 0))
	err = qrcode.NewQrEncode([]byte(code), buf)
	if err != nil {
		log.Sugar.Error("Failed to generate QR code image. ", err)
	}
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="clear-qr-code.jpg"`,
	}
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "image/jpeg", buf, extraHeaders)

	// 使用文件实现
	// c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))
	// c.Writer.Header().Add("Content-Type", "application/octet-stream")
	// c.File(AssetPath+"qr-code.jpg")
}

// licenseApi
func LicenseAPI(c *gin.Context) {
	var (
		cipher string
		err    error
	)

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		switch c.Request.Method {
		// 授权信息
		case "GET":
			c.JSON(http.StatusOK, gin.H{
				"code": statusOk,
				"data": LoadLic().Format(),
				"msg":  "success",
			})
		// 导入授权
		case "POST":
			cipher = c.PostForm("key")
			lic, ok, msg := ChkLicense(cipher)
			if !ok {
				c.JSON(statusOk, gin.H{"code": statusChkLicenseErr, "msg": msg, "data": ""})
				return
			}
			err = PutLicense(cipher)
			if err != nil {
				c.JSON(statusOk, gin.H{"code": statusPutLicenseErr, "msg": err, "data": ""})
				return
			}
			StoreLic(lic)
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": msg, "data": ""})
		case "DELETE":
			code, err := GenClearLicense()
			if err != nil {
				c.JSON(statusOk, gin.H{"code": statusClearLicenseErr, "msg": err, "data": map[string]string{"key": code}})
				return
			}
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": "", "data": map[string]string{"key": code}})
		default:
			c.JSON(statusOk, gin.H{"code": statusMethodErr, "msg": "method error", "data": ""})
		}
	}
}

// 运行状态Api
func NodeStatusAPI(c *gin.Context) {
	var (
		nodeL = make([]stat, 0)
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		nodeM := GetAllNode()
		for _, n := range nodeM {
			nodeL = append(nodeL, stat{n.Attrs.Name, n.Attrs.Addr, utils.RunTime(n.Attrs.Now, n.Attrs.Start)})
		}
		sort.Slice(nodeL, func(i, j int) bool {
			return nodeL[i].ID < nodeL[j].ID
		})
		c.JSON(http.StatusOK, gin.H{"code": statusOk, "data": nodeL, "msg": "success",})
	}
}

// 解绑app
func UntiedAppApi(c *gin.Context) {
	var (
		err error
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		// 获取应用,id,解绑码
		app := c.Param("app")
		id := c.Param("id")
		code := c.PostForm("code")
		if err = Untied(app, id, code); err != nil {
			log.Sugar.Error(err)
			c.JSON(http.StatusOK, gin.H{"code": statusUntiedAppErr, "data": "", "msg": "解绑失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": statusOk, "data": "", "msg": "解绑成功"})
		return
	}
}

// 客户端配置Api
func ConfAPI(c *gin.Context) {
	var (
		list       = make([]conf, 0)
		all        map[string]string
		err        error
		name, text string
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		name = c.Param("name")
		name = strings.Trim(name, "/")
		switch c.Request.Method {
		case "GET":
			// 获取单个config
			if name != "" {
				text, err = GetConfig(name)
				if err != nil {
					c.JSON(statusOk, gin.H{"code": statusGetLicenseErr, "msg": "Code value error." + err.Error(), "data": ""})
					return
				}
				list = append(list, conf{name, text})
				c.JSON(statusOk, gin.H{"code": statusOk, "data": list, "msg": "success"})
				return
			}
			// 获取多个config
			if all, err = GetAllConfig(); err != nil {
				all = map[string]string{"default": err.Error()}
			}
			for name, text := range all {
				list = append(list, conf{name, text})
			}
			sort.Slice(list, func(i, j int) bool {
				return list[i].Name < list[j].Name
			})
			c.JSON(statusOk, gin.H{"code": statusOk, "data": list, "msg": "success"})

		case "DELETE":
			if err = DelConfig(name); err != nil {
				c.JSON(statusOk, gin.H{"code": statusOk, "msg": err, "data": ""})
				return
			}
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": "Delete key success.", "data": ""})

		case "POST", "PUT":
			text = c.PostForm("text")
			if err = PutConfig(name, text); err != nil {
				c.JSON(statusOk, gin.H{"code": statusOk, "msg": err, "data": "",})
				return
			}
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": "Post or Put key success.", "data": ""})
		default:
			c.JSON(statusOk, gin.H{"code": statusMethodErr, "msg": "method error", "data": ""})
		}
	}
}

// 客户端在线Api
func CliOnlineAPI(c *gin.Context) {
	var (
		lines = make([]online, 0)
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		app := c.Param("app")
		app = strings.Trim(app, "/")
		tmp := make(map[string]string)
		if cliMap, err := GetAllClient(app); err == nil {
			for id, status := range cliMap {
				cli := new(Cli)
				if err := json.Unmarshal([]byte(status), cli); err == nil {
					tmp[id] = fmt.Sprintf("节点:%s %s %s", cli.ID, cli.App, utils.RunTime(time.Now().Unix(), cli.Start))
				}
			}
			for id, info := range tmp {
				lines = append(lines, online{id, info})
			}
			sort.Slice(lines, func(i, j int) bool {
				return lines[i].ID < lines[j].ID
			})
			c.JSON(http.StatusOK, gin.H{"code": statusOk, "data": lines, "msg": "success"})
		}
	}
}

// web登录Api
func LoginAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(statusOk, gin.H{"cookie": base64.StdEncoding.EncodeToString([]byte("admin"))})
	}
}

func HelpAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(statusOk, gin.H{"code": statusOk, "data": []string{
			"获取二维码图片：",
			"curl -X GET http://localhost:9527/odin/api/v1/server/qr-code",
			"其他问题请联系技术人员。",
		}, "msg": "success"})
	}
}

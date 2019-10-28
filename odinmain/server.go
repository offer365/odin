package odinmain

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/controller"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/logic"
	pb "github.com/offer365/odin/proto"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
	"time"
)

var gs *grpc.Server

func Run(addr string) {
	var err error
	gs, err = pb.NodeGRpcServer()
	if err != nil {
		log.Sugar.Fatal(err)
		return
	}
	pb.RegisterStaterServer(gs, pb.Self)
	pb.RegisterAuthorizeServer(gs, logic.Author)
	ws := ginServer()
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gs.ServeHTTP(w, r) // grpc server
		} else {
			ws.ServeHTTP(w, r) // gin web server
		}
		return
	})
	listener, err := NewTlsListen([]byte(pb.Server_crt), []byte(pb.Server_key), []byte(pb.Ca_crt),addr)
	if err != nil {
		log.Sugar.Fatal(err)
		return
	}
	err = http.Serve(listener, handle)
	if err != nil {
		log.Sugar.Fatal(err)
		return
	}
}

func NewTlsListen(crt, key, ca []byte, addr string) (net.Listener, error) {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Sugar.Fatal(err)
		return nil, err
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = errors.New("failed to append ca certs")
		log.Sugar.Fatal(err)
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.NoClientCert, // NOTE: 这是可选的!
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Rand:               rand.Reader,
		Time:               time.Now,
		NextProtos:         []string{"http/1.1", http2.NextProtoTLS},
	}
	return tls.Listen("tcp", addr, tlsConfig)
}

// gin 路由
func ginServer() http.Handler {
	gin.SetMode(gin.ReleaseMode) // 生产模式
	r := gin.New()
	r.Use(gin.Recovery()) //Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	//store := cookie.NewStore([]byte("secret"))
	//store.Options(sessions.Options{MaxAge:0})
	//r.Use(sessions.Sessions("odin",store))
	r.LoadHTMLGlob(_assetPath + "html/*")
	r.Static("/static", _assetPath+"static")
	r.StaticFile("/favicon.ico", _assetPath+"static/favicon.ico")
	// api 路由组 basicauth 认证
	api := r.Group("/odin/api/v1", gin.BasicAuth(gin.Accounts{
		User: config.Cfg.Pwd,
	}))

	// 序列号
	api.Any("/server/code", controller.SerialNumAPI)
	// 序列号二维码
	api.GET("/server/qr-code", controller.QrCodeAPI)
	// 授权码
	api.Any("/server/license", controller.LicenseAPI)
	// 注销二维码
	api.GET("/server/qr-license", controller.QrLicenseAPI)
	// 节点状态
	api.GET("/server/nodes", controller.NodeStatusAPI)
	// 解绑接口
	api.POST("/server/untied/:app/:id",controller.UntiedAppApi)
	// 配置接口 kv 存储
	api.Any("/client/conf/*name", controller.ConfAPI)
	// 客户端接口
	api.Any("/client/auth", controller.ClientAPI)
	// 客户端在线接口
	api.GET("/client/online/*app", controller.CliOnlineAPI)
	// web 交互
	api.POST("/web/login", controller.LoginAPI)
	// 帮助
	api.GET("/server/help",controller.HelpAPI)


	r.Any("", func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		r.HandleContext(c)
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "首页",
		})
	})

	return r
}

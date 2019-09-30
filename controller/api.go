package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/offer365/endecrypt"
	"github.com/offer365/endecrypt/endeaesrsa"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/logic"
	"github.com/offer365/odin/model"
	"github.com/offer365/odin/pkg/qrcode"
	"github.com/offer365/odin/pkg/tools"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	success   int = 200
	methodErr int = iota + 410
	notExistErr
	existErr
	getErr
	excessErr
	cipherErr
	saveErr
	bodyErr
	uuidErr
	leaseErr
	renErr
	delErr
)

var secrets = gin.H{"admin": nil}
var Lock sync.Mutex

type body struct {
	Uid   string `json:"uid"`
	Lease int64  `json:"lease"`
	Auth  string `json:"auth"`
}

type result struct {
	Auth   string `json:"auth"`
	Lease  int64  `json:"lease"`
	Cipher string `json:"cipher"`
}

type online struct {
	ID   string `json:"id"`
	Info string `json:"info"`
}

type status struct {
	ID     string `json:"id"`
	Online string `json:"online"`
}

type conf struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// 序列号Api
func SerialNumAPI(c *gin.Context) {
	var (
		key string
		err error
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	key = "Auth error."
	if _, ok := secrets[user]; ok {
		switch c.Request.Method {
		case "GET":
			key, err = logic.GetSerialNum()
		case "POST":
			key, err = logic.ResetSerialNum()
		default:
			key = "Method error."
		}
	}
	if err != nil {
		key = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]string{"key": key, "date": time.Now().Format("2006-01-02 15:04:05")},
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
	code = "Auth error."
	if _, ok := secrets[user]; ok {
		code, err = logic.GetSerialNum()
	}
	if err != nil {
		code = err.Error()
	}

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))  //对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Type", "image/jpeg") // 不能去掉
	buf := bytes.NewBuffer(make([]byte, 0))
	err = qrcode.NewQrEncode([]byte(code), buf)
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="qr-code.jpg"`,
	}
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "image/jpeg", buf, extraHeaders)

	// 使用文件实现
	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	//c.File(AssetPath+"qr-code.jpg")
}

// 注销二维码Api
func QrLicenseAPI(c *gin.Context) {
	var (
		code string
		err  error
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	code = "Auth error."
	if _, ok := secrets[user]; ok {
		code, err = logic.GetClearLicense()
	}
	if err != nil {
		code = err.Error()
	}

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))  //对下载的文件重命名
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
	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "qr-code.jpg"))
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	//c.File(AssetPath+"qr-code.jpg")
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
				"code": 200,
				"data": logic.LoadLic().Format(),
				"msg":  "success",
			})
		// 导入授权
		case "POST":
			cipher = c.PostForm("key")
			lic, ok, msg := logic.ChkLicense(cipher)
			if !ok {
				c.JSON(200, gin.H{"code": 410, "msg": msg, "data": ""})
				return
			}
			err = logic.PutLicense(cipher)
			if err != nil {
				c.JSON(200, gin.H{"code": 510, "msg": err, "data": ""})
				return
			}
			logic.StoreLic(lic)
			c.JSON(200, gin.H{"code": 200, "msg": msg, "data": ""})
		case "DELETE":
			code, err := logic.GenClearLicense()
			if err != nil {
				c.JSON(200, gin.H{"code": 511, "msg": err, "data": map[string]string{"key": code}})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "", "data": map[string]string{"key": code}})
		default:
			c.JSON(200, gin.H{"code": 411, "msg": "Method error.", "data": ""})
		}
	}
}

// 运行状态Api
func NodeStatusAPI(c *gin.Context) {
	var (
		nodeL = make([]status, 0)
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		nodeM := logic.GetAllNode()
		for _, n := range nodeM {
			nodeL = append(nodeL, status{n.Name, fmt.Sprintf("节点:%s ip:%s %s", n.Name, n.IP, tools.RunTime(n.Now, n.Start))})
		}
		sort.Slice(nodeL, func(i, j int) bool {
			return nodeL[i].ID < nodeL[j].ID
		})
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": nodeL, "msg": "success",})
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
				text, err = logic.GetConfig(name)
				if err != nil {
					c.JSON(200, gin.H{"code": 412, "msg": "Code value error." + err.Error(), "data": ""})
					return
				}
				list = append(list, conf{name, text})
				c.JSON(200, gin.H{"code": 200, "data": list, "msg": "success"})
				return
			}
			// 获取多个config
			if all, err = logic.GetAllConfig(); err != nil {
				all = map[string]string{"default": err.Error()}
			}
			for name, text := range all {
				list = append(list, conf{name, text})
			}
			sort.Slice(list, func(i, j int) bool {
				return list[i].Name < list[j].Name
			})
			c.JSON(200, gin.H{"code": 200, "data": list, "msg": "success"})

		case "DELETE":
			if err = logic.DelConfig(name); err != nil {
				c.JSON(200, gin.H{"code": 200, "msg": err, "data": ""})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "Delete key success.", "data": ""})

		case "POST", "PUT":
			_, ok := logic.PutWhiteList[name]
			if ok {
				c.JSON(200, gin.H{"code": 200, "msg": "The key " + name + " can only be accessed and cannot be edited.", "data": ""})
				return
			}
			text = c.PostForm("text")
			if err = logic.PutConfig(name, text); err != nil {
				c.JSON(200, gin.H{"code": 200, "msg": err, "data": "",})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "Post or Put key success.", "data": ""})
		default:
			c.JSON(200, gin.H{"code": 1, "msg": "Method error.", "data": ""})
		}
	}
}

// 客户端Api
func ClientAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		app := c.Param("app")
		id := c.Param("id")
		key := app + "/" + id
		// 判断app 是否在授权中 && 是否到期
		if !logic.LoadLic().CheckTime(app) {
			c.JSON(200, gin.H{"code": notExistErr, "data": result{Lease: 0}, "msg": "APP does not exist or authorization expires."})
			return
		}
		// 判断当前的客户端是否已经存在
		cli, exist := logic.GetClient(key)
		switch c.Request.Method {
		case "POST": // 认证
			//TODO 原子操作代替锁
			Lock.Lock()
			defer Lock.Unlock()
			if cli != nil || exist {
				c.JSON(200, gin.H{"code": existErr, "data": result{Lease: 0}, "msg": "The id app already exists."})
				return
			}

			nc := new(model.Cli)
			nc.App = app
			nc.IP = c.ClientIP()
			nc.ID = id
			nc.Start = time.Now().Unix()
			nc.Uuid = uuid.Must(uuid.NewV4()).String()
			// 获取cli的实例个数
			num, err := logic.ClientCount(app)
			if err != nil {
				c.JSON(200, gin.H{"code": getErr, "data": result{Lease: 0}, "msg": "Failed to get the number of app instances." + err.Error()})
				return
			}
			// 检查实例是否超出授权个数
			if int(num) >= logic.LoadLic().APPs[app].Instance {
				c.JSON(200, gin.H{"code": excessErr, "data": result{Lease: 0}, "msg": "The app has insufficient remaining instances."})
				return
			}

			// 生成密文
			cipher, err := endeaesrsa.PriEncrypt([]byte(nc.Uuid), endecrypt.PirkeyClient2048, endecrypt.AesKeyClient2)
			if err != nil {
				c.JSON(200, gin.H{"code": cipherErr, "data": result{Lease: 0}, "msg": "The app failed to generate cipher." + err.Error()})
				return
			}
			// 这个租约id 在PutClient里面赋值
			lease, err := logic.PutClient(key, nc)
			if err != nil {
				c.JSON(200, gin.H{"code": saveErr, "data": result{Lease: 0}, "msg": "This app failed to save the instance." + err.Error()})
				return
			}
			attr := logic.LoadLic().APPs[app].Attr
			// 没有意义，混淆作用
			attr["time"] = time.Now().UnixNano()
			byt, _ := json.Marshal(attr)
			auth, err := endeaesrsa.PriEncrypt(byt, endecrypt.PirkeyClient2048, endecrypt.AesKeyClient2)
			if err != nil {
				c.JSON(200, gin.H{"code": cipherErr, "data": result{Lease: 0}, "msg": "The app failed to generate cipher." + err.Error()})
				return
			}
			// 生成的auth 与 cipher 可以使用不同的加密算法。
			c.JSON(200, gin.H{"code": success, "data": result{auth, lease, cipher}, "msg": "success",})
			return
		case "PUT": // 心跳
			// 实例不存在
			if !exist {
				c.JSON(200, gin.H{"code": notExistErr, "data": result{Lease: 0}, "msg": "The client does not exist."})
				return
			}
			data, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": bodyErr, "data": result{Lease: 0}, "msg": "The request bd data error." + err.Error()})
				return
			}
			bd := new(body)
			err = json.Unmarshal(data, bd)
			if err != nil {
				c.JSON(200, gin.H{"code": bodyErr, "data": result{Lease: 0}, "msg": "The request bd data error." + err.Error()})
				return
			}
			// uuid不匹配
			if bd.Uid != cli.Uuid {
				c.JSON(200, gin.H{"code": uuidErr, "data": result{Lease: 0}, "msg": "The client uuid error."})
				return
			}
			// 租约id不匹配
			if bd.Lease != cli.Lease {
				c.JSON(200, gin.H{"code": leaseErr, "data": result{Lease: 0}, "msg": "The client lease error."})
				return
			}
			// 续租失败
			if err = logic.KeepAliveClient(key, bd.Lease); err != nil {
				c.JSON(200, gin.H{"code": renErr, "data": result{Lease: bd.Lease}, "msg": "Renewal failed." + err.Error()})
				return
			}
			// 自定义 auth 返回重要的信息等等
			c.JSON(200, gin.H{"code": success, "data": result{"", bd.Lease, ""}, "msg": "Successful renewal."})
			return
		case "DELETE": // 关闭
			if !exist {
				c.JSON(200, gin.H{"code": notExistErr, "data": result{Lease: 0}, "msg": "The client does not exist."})
				return
			}
			data, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": bodyErr, "data": result{Lease: 0}, "msg": "The request bd data error." + err.Error()})
				return
			}
			bd := new(body)
			err = json.Unmarshal(data, bd)
			if err != nil {
				c.JSON(200, gin.H{"code": bodyErr, "data": result{Lease: 0}, "msg": "The request bd data error." + err.Error()})
				return
			}
			// uuid不匹配
			if bd.Uid != cli.Uuid {
				c.JSON(200, gin.H{"code": uuidErr, "data": result{Lease: 0}, "msg": "The client uuid error."})
				return
			}
			// 租约id不匹配
			if bd.Lease != cli.Lease {
				c.JSON(200, gin.H{"code": leaseErr, "data": result{Lease: 0}, "msg": "The client lease error."})
				return
			}
			// 删掉此客户端实例
			if err := logic.DelClient(key, bd.Lease); err != nil {
				c.JSON(200, gin.H{"code": delErr, "data": result{Lease: bd.Lease}, "msg": "Deleting an instance failed." + err.Error()})
				return
			}
			c.JSON(200, gin.H{"code": 200, "data": result{Lease: bd.Lease}, "msg": "Deleting an instance succeed."})
			return
		default:
			c.JSON(200, gin.H{"code": methodErr, "data": result{Lease: 0}, "msg": "Method error."})
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
		if cliMap, err := logic.GetAllClient(app); err == nil {
			for id, status := range cliMap {
				cli := new(model.Cli)
				if err := json.Unmarshal([]byte(status), cli); err == nil {
					tmp[id] = fmt.Sprintf("节点:%s(%s) %s %s", cli.ID, cli.IP, cli.App, tools.RunTime(time.Now().Unix(), cli.Start))
				}
			}
			for id, info := range tmp {
				lines = append(lines, online{id, info})
			}
			sort.Slice(lines, func(i, j int) bool {
				return lines[i].ID < lines[j].ID
			})
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": lines, "msg": "success"})
		}
	}
}

// web登录Api
func LoginAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(200, gin.H{"cookie": base64.StdEncoding.EncodeToString([]byte("admin"))})
	}
}

func HelpAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(200, gin.H{"code": 200, "data": []string{
			"获取二维码图片：",
			"curl -X GET http://127.0.0.1:9999/odin/api/v1/server/qr-code",
			"其他问题请联系技术人员。",
		}, "msg": "success"})
	}
}

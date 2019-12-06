package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/offer365/example/qrcode"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/logic"
	"github.com/offer365/odin/proto"
	"github.com/offer365/odin/utils"

	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	statusOk            int   = 200
	statusMethodErr     int32 = 405
	statusChkLicenseErr int32 = iota + 440
	statusPutLicenseErr
	statusClearLicenseErr
	statusUntiedAppErr
	statusGetLicenseErr
)

var secrets = gin.H{"admin": nil}
var Lock sync.Mutex

type online struct {
	ID   string `json:"id"`
	Info string `json:"info"`
}

type status struct {
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
			code, err = logic.GetSerialNum()
		case "POST":
			code, err = logic.ResetSerialNum()
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
		code, err = logic.GetSerialNum()
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
		code, err = logic.GetClearLicense()
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
				"data": logic.LoadLic().Format(),
				"msg":  "success",
			})
		// 导入授权
		case "POST":
			cipher = c.PostForm("key")
			lic, ok, msg := logic.ChkLicense(cipher)
			if !ok {
				c.JSON(statusOk, gin.H{"code": statusChkLicenseErr, "msg": msg, "data": ""})
				return
			}
			err = logic.PutLicense(cipher)
			if err != nil {
				c.JSON(statusOk, gin.H{"code": statusPutLicenseErr, "msg": err, "data": ""})
				return
			}
			logic.StoreLic(lic)
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": msg, "data": ""})
		case "DELETE":
			code, err := logic.GenClearLicense()
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
		nodeL = make([]status, 0)
	)
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		nodeM := logic.GetAllNode()
		for _, n := range nodeM {
			nodeL = append(nodeL, status{n.Attrs.Name, n.Attrs.Addr, utils.RunTime(n.Attrs.Now, n.Attrs.Start)})
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
		if err = logic.Untied(app, id, code); err != nil {
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
				text, err = logic.GetConfig(name)
				if err != nil {
					c.JSON(statusOk, gin.H{"code": statusGetLicenseErr, "msg": "Code value error." + err.Error(), "data": ""})
					return
				}
				list = append(list, conf{name, text})
				c.JSON(statusOk, gin.H{"code": statusOk, "data": list, "msg": "success"})
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
			c.JSON(statusOk, gin.H{"code": statusOk, "data": list, "msg": "success"})

		case "DELETE":
			if err = logic.DelConfig(name); err != nil {
				c.JSON(statusOk, gin.H{"code": statusOk, "msg": err, "data": ""})
				return
			}
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": "Delete key success.", "data": ""})

		case "POST", "PUT":
			text = c.PostForm("text")
			if err = logic.PutConfig(name, text); err != nil {
				c.JSON(statusOk, gin.H{"code": statusOk, "msg": err, "data": "",})
				return
			}
			c.JSON(statusOk, gin.H{"code": statusOk, "msg": "Post or Put key success.", "data": ""})
		default:
			c.JSON(statusOk, gin.H{"code": statusMethodErr, "msg": "method error", "data": ""})
		}
	}
}

// 客户端Api
func ClientAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		// app := c.Param("app")
		// id := c.Param("id")
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Sugar.Error("ready request body error: ", err)
			return
		}
		req := new(proto.Request)
		err = json.Unmarshal(data, req)
		if err != nil {
			return
		}
		switch c.Request.Method {
		case "POST": // 认证
			resp, err := logic.Author.Auth(context.TODO(), req)
			if err != nil {
				log.Sugar.Error(resp, err)
				return
			}
			c.JSON(statusOk, resp)
			return
		case "PUT": // 心跳
			resp, err := logic.Author.KeepLine(context.TODO(), req)
			if err != nil {
				log.Sugar.Error(resp, err)
				return
			}
			c.JSON(statusOk, resp)
			return
		case "DELETE": // 关闭
			resp, err := logic.Author.OffLine(context.TODO(), req)
			if err != nil {
				log.Sugar.Error(resp, err)
				return
			}
			c.JSON(statusOk, resp)
			return
		default:
			c.JSON(statusOk, &proto.Response{Code: statusMethodErr, Data: nil, Msg: "method error",})
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
				cli := new(logic.Cli)
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

package odinmain

import (
	"crypto/rand"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/controller"
	"github.com/offer365/odin/log"
	"net/http"
	"time"
)

const (
	cert_pem = `
-----BEGIN CERTIFICATE-----
MIIENjCCAp6gAwIBAgIQYJNuLVi6LLxIdVGz1tGqazANBgkqhkiG9w0BAQsFADB1
MR4wHAYDVQQKExVta2NlcnQgZGV2ZWxvcG1lbnQgQ0ExJTAjBgNVBAsMHHJvb3RA
aXoyemU3cGxuNWlpZnp3dW1lenAyNXoxLDAqBgNVBAMMI21rY2VydCByb290QGl6
MnplN3BsbjVpaWZ6d3VtZXpwMjV6MB4XDTE5MDYwMTAwMDAwMFoXDTI5MDcxMTEz
MjMwMlowUDEnMCUGA1UEChMebWtjZXJ0IGRldmVsb3BtZW50IGNlcnRpZmljYXRl
MSUwIwYDVQQLDBxyb290QGl6MnplN3BsbjVpaWZ6d3VtZXpwMjV6MIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA22r6PBSOeQ8kxHJSiRq4fd7Er6D4I7ME
352QVkg0Y/grsluZ6CTKzKOq8KyWfEF2fl/qHSu3TpAALtqfpXY+bni5LLHpu87Z
ShTy2sZe08U4SUrdCKFqdR6ZJBczaHpwuHD39h7cpkjEMRB4jIXK89d1+07e02R6
V396+628O7qBVJUSu9kNxgHNXGNtADOLUGz95maR2aDoDFL+GBVTu+Ww0xd0faao
8L5zRKvV1jY8eV15qFQ4VJxdnzl0Z8hY8wG2JtyeeEE3GB+QuJs05IyWPqGAEFhu
J7JKmU2gCcI9lKA3voAF3/WeQsiTzvyDatjv3RK6Gr6Vi2/CStJ2NQIDAQABo2cw
ZTAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/
BAIwADAfBgNVHSMEGDAWgBS1KZ6/gmcwY8Cku8BV34whHAj0izAPBgNVHREECDAG
ggRvZGluMA0GCSqGSIb3DQEBCwUAA4IBgQCQwkUN4/mpRSzIZ7y2+1qli3tSGxPR
MG1pI62HngJfRXDAnKKsJ95YBB7+HOqWLFNjpA8g/aAOCLFILeDBdhOK73Kl3bxK
FvPeC4iyw9kTlHtJ9r4nMMsWOE4nC/ioS3MRPbwlDnmUjqfQ6gUVbjfV/tc4eYtn
66a0v3TqfjWB1EEw6NbDK7h6kL9VBmB/G33snGJmKJWDGxJPSafIRzjvAnGR9W2M
qkT65QrNmvDWd1JzXQcW58sH7wJBkSXyOYkID49fPIry/4h+nWLcGpw5VuqClabC
K6zrlGVD5tw3eq84XVpjAon5syMkKgcWv6ekljB4+mHDuky+d/em73nPVBthjBuo
aVYSj0R0ogu34gUuCDkZcr4ZiE7LHK/xcUtcARo3Fjig+hqq7wj9ZfrJ3bVphDvC
8+Bu7hjdZalZC+LgwLpDwMgbSzDpCRoVZiRk3TNMjCEJt+ISLUjmaa4bnORrkhfo
VlgRkNF5acH0YbFlacumtuB5uBEJrGvDJzM=
-----END CERTIFICATE-----
`
	key_pem = `
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDbavo8FI55DyTE
clKJGrh93sSvoPgjswTfnZBWSDRj+CuyW5noJMrMo6rwrJZ8QXZ+X+odK7dOkAAu
2p+ldj5ueLkssem7ztlKFPLaxl7TxThJSt0IoWp1HpkkFzNoenC4cPf2HtymSMQx
EHiMhcrz13X7Tt7TZHpXf3r7rbw7uoFUlRK72Q3GAc1cY20AM4tQbP3mZpHZoOgM
Uv4YFVO75bDTF3R9pqjwvnNEq9XWNjx5XXmoVDhUnF2fOXRnyFjzAbYm3J54QTcY
H5C4mzTkjJY+oYAQWG4nskqZTaAJwj2UoDe+gAXf9Z5CyJPO/INq2O/dEroavpWL
b8JK0nY1AgMBAAECggEBAM8e5KfiH7tW+DYYVKDngFARAUlogdPxISCU87L+5bWY
hmcO4PGqCWWy+aHGySbyBJC2qaBvq9GVTRbteNYQEE7n1qTCLQkD8UllDPpHVyxA
dyl4ab3D4WI9SAIxhG2TZuQ0f1ztNQwilFBcY+8CPNqBAPYBNYYGyxXdWJJLJeya
Gg0KHEfuZv7ueZWWvPqKMW2PGtm+spURze6AVfP1WU9XGdXJrFisLZukM9Tp9Uws
14V6dP0bGHggZ6Lx5ZHGTl8ed/xmtDOkeI8LOeK+cQ/JtfE6TPrB8WFpEdgJKRir
f2EnHpfuIfAQ2hsNpQOvO3Rk3AnhbBeSyiQ+vO4RzgECgYEA6bORkzN6ESdpSSPv
kn/trKxwAc2l84hM8cJ0BtXEVNed8y0MkZwHQk38DL6kU7p/e0inaVuLf1tbBPQb
LjpzJc9Fmit93P25TvxC17S8u2RJDZYZ0v1FVnzTEzI9ZLhTJZ9AVB8Jf0OlA5Fx
AsiMDNoo4G0rbnxabfn6gXQSjcUCgYEA8FqEI2S0u5zRjdZUD+1S5YcE7aBX9L0q
XrM3xI7GZZjELVt5z9YZmR8YAwntpOjBsnLyI4kiSAGpCCS9PdpBs0yZa0oQ24XA
KuZh1WT22269Zvp/JHZ+M+Rdboz9OY5T8ofero9JyxyvHBVi0mbbMmMiVieDfhl6
0wW3/ONEvbECgYEAl4C34QvAGJrKIIZRa1HPzN9FBYZCDSzRZPFAsqWmT7IwTVNp
EIRsGEniGokEktsWhd/F2AFm37tjuERf0opF178VSirjv34kwdW7p4cdywXqbgpe
128lojntxEYPktoD3SHuXBp616wMr9F7x+gnErXjRgq/2zJ2lVE3WvDajlkCgYAa
Pof3JWPmqHTpO+Hp60wF6/xJxhxUiOM7e+429DANn+Sr3zUp0ILzCUYh7s+YFiIw
TgTKhIrNugCu9vQC8PYDkfWelXPJxIz7IjTEjEW4KBteRzPi011sZR8elx5/Tl80
OEnEXbj9CKDGPD+SIdEFa3WwWpgtCLM0n4c7gcVbwQKBgQDHxDqOC1bo6al0ZZtL
XXppw1Lf7VD5i0gLS/qNe6gIs8DJx9WF02JD13Yi1gyqDosErljcU3RA5hZsqoc3
iGMyEfmKOR4vYGq73R0XVqUqKgpBG5xQs3StDpsd1n7Ye81vq1YhVKu5LrjiTRkJ
tNBT2he8EJWiBzyZ31nGwTt/pQ==
-----END PRIVATE KEY-----
`
)

///////////////////////////////////////////////////////////////////////////
/// web																///
///////////////////////////////////////////////////////////////////////////

// 启动 gin web
func RunWebWithHttp(port string) {
	if err := route().Run(":" + port); err != nil {
		log.Sugar.Fatal("starting gin web service failed. error: ", err.Error())
	}
}

func RunWebWithHttps(port string) {
	crt, err := tls.X509KeyPair([]byte(cert_pem), []byte(key_pem))
	if err != nil {
		log.Sugar.Fatal(err.Error())
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, TLS uses time.Now.
	tlsConfig.Time = time.Now
	// Rand provides the source of entropy for nonces and RSA blinding.
	// If Rand is nil, TLS uses the cryptographic random reader in package
	// crypto/rand.
	// The Reader must be safe for use by multiple goroutines.
	tlsConfig.Rand = rand.Reader
	l, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		log.Sugar.Fatal(err.Error())
	}
	err = http.Serve(l, route())
	if err != nil {
		log.Sugar.Fatal(err.Error())
	}
}

// gin 路由
func route() (r *gin.Engine) {
	// debug
	if debug {
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode) // 生产模式
		r = gin.New()
		r.Use(gin.Recovery()) //Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	}
	r.LoadHTMLGlob(AssetPath + "html/*")

	// api 路由组
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
	// 配置接口
	//api.Any("/client/conf", ConfAPI)
	api.Any("/client/conf/*name", controller.ConfAPI)
	// 客户端接口
	api.Any("/client/auth/:app/:id", controller.ClientAPI)
	// 客户端在线接口
	api.GET("/client/online/*app", controller.CliOnlineAPI)
	// 客户端在线接口
	api.POST("/web/login", controller.LoginAPI)
	api.GET("/server/help", controller.HelpAPI)

	//r.Use(SimpleSession)
	r.Static("/static", AssetPath+"static")
	r.Any("", func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		r.HandleContext(c)
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "首页",
		})
	})

	r.StaticFile("/favicon.ico", AssetPath+"static/favicon.ico")
	return
}

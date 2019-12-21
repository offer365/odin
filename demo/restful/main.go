package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/httplib"
	"github.com/offer365/example/endecrypt/endersa"
	pb "github.com/offer365/odin/demo/odinX"
	"github.com/offer365/odin/utils"

	"strconv"
	"time"
)

var (
	Cfg       *Config
	tlsConifg *tls.Config
)

// 加解密方法
type CryptFunc func(src []byte) ([]byte, error)

// hash 方法
type HashFunc func(byt []byte) string

type Config struct {
	AppName    string
	ServerAddr []string
	ServerName string
	GRpcUser   string
	GRpcPwd    string
	ClientCrt  string
	ClientKey  string
	CaCrt      string

	// odin & client app
	VerifyEncrypt CryptFunc // token 密文加密
	CipherDecrypt CryptFunc // uuid 解密
	AuthDecrypt   CryptFunc // auth 数据解密
	UuidHash      HashFunc
}

// odin & app

func verifyEncrypt1(src []byte) ([]byte, error) {
	return endersa.PubEncrypt(src, []byte(_rsa2048pub1))
}

func verifyDecrypt1(src []byte) ([]byte, error) {
	return endersa.PriDecrypt(src, []byte(_rsa2048pri1))
}

func cipherEncrypt1(src []byte) ([]byte, error) {
	return endersa.PubEncrypt(src, []byte(_rsa2048pub2))
}

func cipherDecrypt1(src []byte) ([]byte, error) {
	return endersa.PriDecrypt(src, []byte(_rsa2048pri2))
}

func authEncrypt1(src []byte) ([]byte, error) {
	return endersa.PubEncrypt(src, []byte(_rsa2048pub3))
}

func authDecrypt1(src []byte) ([]byte, error) {
	return endersa.PriDecrypt(src, []byte(_rsa2048pri3))
}

func HashFunc2(src []byte) string {
	return utils.Sha256Hex(src, []byte(storeHashSalt2))
}

func main() {
	Cfg = &Config{
		AppName:       "nlp",
		ServerAddr:    []string{"https://" + "127.0.0.1" + ":1443/odin/api/v1/client/auth"},
		ServerName:    server_name,
		GRpcUser:      grpcUser,
		GRpcPwd:       grpcPwd,
		ClientCrt:     client_crt,
		ClientKey:     client_key,
		CaCrt:         ca_crt,
		VerifyEncrypt: verifyEncrypt1,
		CipherDecrypt: cipherDecrypt1,
		AuthDecrypt:   authDecrypt1,
		UuidHash:      HashFunc2,
	}
	tlsConifg = Tls([]byte(Cfg.ClientCrt), []byte(Cfg.ClientKey), []byte(Cfg.CaCrt))
	ManyApp(1000)
	select {}
}

func NewApp(name, id, token string, servers []string) (app *Application) {
	app = &Application{
		Name:    name,
		ID:      id,
		Servers: servers,
		Token:   token,
		Cycle:   new(Cycle),
	}
	return
}

type Application struct {
	Name    string
	ID      string
	Servers []string
	Token   string
	*Cycle
	tlsConfig *tls.Config
}

// 一次认证请求
type Cycle struct {
	Urls     []string
	Uuid     string
	AuthInfo string
	Lease    int64
}

func Tls(crt, key, ca []byte) *tls.Config {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		// err = errors.New("failed to append ca certs")
		return nil
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.RequireAndVerifyClientCert, // NOTE: 这是可选的!
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Rand:               rand.Reader,
		Time:               time.Now,
		// NextProtos:         []string{"http/1.1", http2.NextProtoTLS},
	}

}

// auth
func (app *Application) Active() {
	verify := make(map[string]interface{})
	verify["app"] = app.Name
	verify["id"] = app.ID
	verify["date"] = time.Now().Unix()
	verify["token"] = app.Token
	byt, err := json.Marshal(verify)
	if err != nil {
		fmt.Println(err)
		return
	}
	byt, err = Cfg.VerifyEncrypt(byt)
	if err != nil {
		fmt.Println(err)
		return
	}
	req := &pb.Request{
		App:    app.Name,
		Id:     app.ID,
		Date:   time.Now().Unix(),
		Verify: base64.StdEncoding.EncodeToString(byt),
	}
	result, err := httplib.Post(app.Servers[0]).SetTimeout(2*time.Second, 3*time.Second).Debug(true).SetBasicAuth(Cfg.GRpcUser, Cfg.GRpcPwd).SetTLSClientConfig(tlsConifg).Header("Content-Type", "application/json; charset=utf-8").JSONBody(req)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	resp := new(pb.Response)
	fmt.Println(result.ToJSON(resp))
	byt, err = Cfg.CipherDecrypt(resp.Data.Cipher) // uuid
	app.Uuid = string(byt)
	byt, err = Cfg.AuthDecrypt(resp.Data.Auth) // {"attrs":[{"Name":"热词","Key":"hotword","Value":1000},{"Name":"类热词","Key":"classword","Value":1000}],"time":1571909203232224000}
	app.AuthInfo = string(byt)                 // eg: {"attrs":[{"Name":"热词","Key":"hotword","Value":1000},{"Name":"类热词","Key":"classword","Value":1000}],"time":1571909203232224000}
	fmt.Println(string(byt), err)
	app.Lease = resp.Data.Lease
}

func (app *Application) KeepLine() {
	req := &pb.Request{
		App:    app.Name,
		Id:     app.ID,
		Date:   time.Now().Unix(),
		Verify: "",
		Umd5:   Cfg.UuidHash([]byte(app.Uuid)),
		Lease:  app.Lease,
	}

	for range time.Tick(time.Second * 6) {
		byt, _ := json.Marshal(req)
		fmt.Println(string(byt))
		result, _ := httplib.Put(app.Servers[0]).SetTimeout(2*time.Second, 3*time.Second).Debug(true).SetBasicAuth(Cfg.GRpcUser, Cfg.GRpcPwd).SetTLSClientConfig(tlsConifg).Header("Content-Type", "application/json; charset=utf-8").JSONBody(req)
		resp := &pb.Response{}
		result.ToJSON(resp)
		fmt.Println(result.String())
		fmt.Println(resp.Data.Lease)
	}
}

func (app *Application) OffLine() {
	req := &pb.Request{
		App:    app.Name,
		Id:     app.ID,
		Date:   time.Now().Unix(),
		Verify: "",
		Umd5:   Cfg.UuidHash([]byte(app.Uuid)),
		Lease:  app.Lease,
	}

	result, _ := httplib.Delete(app.Servers[0]).SetTimeout(2*time.Second, 3*time.Second).Debug(true).SetBasicAuth(Cfg.GRpcUser, Cfg.GRpcPwd).SetTLSClientConfig(tlsConifg).Header("Content-Type", "application/json; charset=utf-8").JSONBody(req)
	resp := &pb.Response{}
	result.ToJSON(resp)
	fmt.Println(resp.Data.Lease)

}

func ManyApp(ins int) {
	apps := make([]*Application, 0)
	for i := 0; i < ins; i++ {
		appID := "app" + strconv.Itoa(i)
		app := NewApp(Cfg.AppName, appID, appID, Cfg.ServerAddr)
		apps = append(apps, app)
	}
	for _, ap := range apps {
		time.Sleep(time.Millisecond * 500)
		go func(app *Application) {
			time.Sleep(time.Second)
			app.Active()
			app.KeepLine()
		}(ap)
	}
}

func SingleAPP() {
	app := NewApp("hotword", "app"+strconv.Itoa(99), "app"+strconv.Itoa(99), Cfg.ServerAddr, )
	app.Active()
	app.KeepLine()
	app.OffLine()
	fmt.Println(111)
}

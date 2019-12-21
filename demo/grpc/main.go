package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/offer365/example/endecrypt/endersa"
	corec "github.com/offer365/example/grpc/core/client"
	pb "github.com/offer365/odin/demo/odinX"
	"github.com/offer365/odin/utils"
	"google.golang.org/grpc"
)

var (
	Cfg *Config
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

type Authentication struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return true
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
		ServerAddr:    []string{"127.0.0.1:1443"},
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
	ManyApp(1000)
	select {}
}

func ManyApp(ins int) {
	apps := make([]*Application, 0)
	for i := 0; i < ins; i++ {
		appID := "app" + strconv.Itoa(i)
		app := NewApp(Cfg.AppName, appID, appID, Cfg.ServerAddr, []byte(Cfg.ClientCrt), []byte(Cfg.ClientKey), []byte(Cfg.CaCrt))
		apps = append(apps, app)
	}
	for _, ap := range apps {
		time.Sleep(time.Millisecond * time.Duration(500))
		go func(app *Application) {
			time.Sleep(time.Second)
			app.Active()
			app.KeepLine()
		}(ap)
	}
}

func SingleAPP() {
	app := NewApp("nlp", "app"+strconv.Itoa(99), "app"+strconv.Itoa(99), Cfg.ServerAddr, []byte(Cfg.ClientCrt), []byte(Cfg.ClientKey), []byte(Cfg.CaCrt))
	app.Active()
	app.KeepLine()
	app.OffLine()
	fmt.Println(111)
}

func NewApp(name, id, token string, servers []string, crt, key, ca []byte) (app *Application) {
	app = &Application{
		Name:    name,
		ID:      id,
		Servers: servers,
		Token:   token,
		Cycle:   new(Cycle),
	}
	auth := &Authentication{
		User:     Cfg.GRpcUser,
		Password: Cfg.GRpcPwd,
	}
	conn, err := corec.NewRpcClient(
		corec.WithAddr(app.Servers[0]),
		corec.WithServerName(Cfg.ServerName),
		corec.WithCert(crt),
		corec.WithKey(key),
		corec.WithCa(ca),
		corec.WithDialOption(grpc.WithPerRPCCredentials(auth)),
	)

	if err != nil {
		return
	}
	app.conn = conn
	app.cli = pb.NewAuthorizeClient(conn)
	return
}

type Application struct {
	conn    *grpc.ClientConn
	cli     pb.AuthorizeClient
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

	resp, err := app.cli.Auth(context.TODO(), req)
	byt, err = Cfg.CipherDecrypt(resp.Data.Cipher)

	app.Uuid = string(byt)
	byt, err = Cfg.AuthDecrypt(resp.Data.Auth)

	app.AuthInfo = string(byt) // eg: {"attrs":[{"Name":"热词","Key":"hotword","Value":1000},{"Name":"类热词","Key":"classword","Value":1000}],"time":1571909203232224000}
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
		resp, err := app.cli.KeepLine(context.TODO(), req)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp.Code, resp.Msg)
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
	resp, err := app.cli.OffLine(context.TODO(), req)
	fmt.Println(resp.Data.Lease, err)

	fmt.Println(resp.Data.Lease)

}

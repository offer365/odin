package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/offer365/example/endecrypt"
	corec "github.com/offer365/example/grpc/core/client"
	pb "github.com/offer365/odin/demo/proto"
	"github.com/offer365/odin/utils"
	"google.golang.org/grpc"
)

var (
	auth      *Authentication
	_username = "C205v406x68f5IM7"
	_password = "c9bJ3v7FQ11681EP"
)

func init() {
	auth = &Authentication{
		User:     _username,
		Password: _password,
	}
}

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	ManyApp(1000)
	select {}
}

func ManyApp(ins int) {
	apps := make([]*Application, 0)
	for i := 0; i < ins; i++ {
		token := "app" + strconv.Itoa(i)
		app := NewApp("nlp", "app"+strconv.Itoa(i), token, []string{"10.0.0.200:9527"}, []byte(pb.Client_crt), []byte(pb.Client_key), []byte(pb.Ca_crt))
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
	app := NewApp("hotword", "app"+strconv.Itoa(99), "app"+strconv.Itoa(99), pb.Member, []byte(pb.Client_crt), []byte(pb.Client_key), []byte(pb.Ca_crt))
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
	conn, err := corec.NewRpcClient(
		corec.WithAddr(app.Servers[0]),
		corec.WithServerName("server.io"),
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
	fmt.Println(string(byt))
	byt, err = endecrypt.Encrypt(endecrypt.Aes2key32, byt)
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
	fmt.Println(err)
	byt, err = endecrypt.Decrypt(endecrypt.Pub2Rsa1024, resp.Data.Cipher)
	app.Uuid = string(byt)
	byt, err = endecrypt.Decrypt(endecrypt.Pub2Rsa2048, resp.Data.Auth)
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
		Umd5:   utils.Md5sum([]byte(app.Uuid), nil),
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
		Umd5:   utils.Md5sum([]byte(app.Uuid), nil),
		Lease:  app.Lease,
	}
	resp, err := app.cli.OffLine(context.TODO(), req)
	fmt.Println(resp.Data.Lease, err)

	fmt.Println(resp.Data.Lease)

}

package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/offer365/endecrypt"
	"github.com/offer365/endecrypt/endeaesrsa"
	"log"
	mr "math/rand"
	"strconv"
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

var member []string = []string{"10.0.0.200", "10.0.0.201", "10.0.0.202"}
var format = "https://%s:8888/odin/api/v1/client/auth/"
var app string

func main() {
	app = "nlp"
	Example(GetCliWithNum(app, 1000)...)
	//app = "test_app"
	//Example(GetCliWithNum(url, app, 50)...)
	//app = "demo3"
	//Example(GetCliWithNum(url, app, 50)...)
	select {}
}

func GetCliWithNum(app string, n int) []*Client {
	cli := make([]*Client, 0)
	for i := 0; i < n; i++ {
		index := i % len(member)
		ip := member[index]
		url := fmt.Sprintf(format, ip)
		cli = append(cli, NewCli(url, app, "a"+strconv.Itoa(i)))
	}
	return cli
}

func Example(cli ...*Client) {
	for _, c := range cli {
		time.Sleep(500 * time.Millisecond)
		go c.RunExample()
	}
}

func NewCli(url, app, id string) *Client {
	return &Client{
		Attr: &Attr{
			Url: url + app + "/" + id,
			App: app,
			ID:  id,
			Uid: "",
		},
		Result: new(Result),
	}
}

type Client struct {
	*Attr
	*Result
	tls *tls.Config
}

type Attr struct {
	Url      string
	App      string
	ID       string
	Uid      string
	AuthInfo string
}

type Result struct {
	Code int `json:"code"`
	Data `json:"data"`
	Msg  string `json:"msg"`
	//{"code": 200, "lease": lease, "msg": str}
}

type Data struct {
	Auth   string `json:"auth"`
	Lease  int64  `json:"lease"`
	Cipher string `json:"cipher"`
}

type Body struct {
	Lease int64
	Uid   string
}

func (cli *Client) Tls() {
	crt, err := tls.X509KeyPair([]byte(cert_pem), []byte(key_pem))
	if err != nil {
		log.Fatalln(err.Error())
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
	tlsConfig.InsecureSkipVerify = true
	cli.tls = tlsConfig
}

func (cli *Client) POST() {
	cli.Result = new(Result)
	result := httplib.Post(cli.Url).SetTimeout(2*time.Second, 3*time.Second).Debug(false).SetBasicAuth("admin", "123").SetTLSClientConfig(cli.tls).Header("Content-Type", "application/json; charset=utf-8")
	result.ToJSON(cli.Result)
	//cli.Body.Lease = cli.result.Data.Lease
	//

	//auth 与 cipher 可以使用不通的加密算法。
	if cli.Cipher != "" {
		cli.Uid, _ = endeaesrsa.PubDecrypt(cli.Cipher, endecrypt.PubkeyClient2048, endecrypt.AesKeyClient2)
	}
	if cli.Auth != "" {
		cli.AuthInfo, _ = endeaesrsa.PubDecrypt(cli.Auth, endecrypt.PubkeyClient2048, endecrypt.AesKeyClient2)
	}

	fmt.Println("认证..", cli.ID, cli.App, cli.Uid, cli.Lease, cli.AuthInfo, cli.Msg, cli.Url)
}

func (cli *Client) PUT() {
	body := &Body{
		Lease: cli.Lease,
		Uid:   cli.Uid,
	}
	res := new(Result)
	byt, _ := json.Marshal(body)
	result := httplib.Put(cli.Url).SetTimeout(2*time.Second, 3*time.Second).Debug(false).SetBasicAuth("admin", "123").SetTLSClientConfig(cli.tls).Body(byt).Header("Content-Type", "application/json; charset=utf-8")
	str, err := result.String()
	result.ToJSON(res)
	if res.Code != 200 {
		fmt.Println("心跳..", cli.ID, cli.App, cli.Uid, cli.Lease, str, err, cli.Url)
	}

}

func (cli *Client) DELETE() {
	body := &Body{
		Lease: cli.Lease,
		Uid:   cli.Uid,
	}
	byt, _ := json.Marshal(body)
	result := httplib.Delete(cli.Url).SetTimeout(2*time.Second, 3*time.Second).Body(byt).Debug(false).SetBasicAuth("admin", "123").SetTLSClientConfig(cli.tls).Body(byt).Header("Content-Type", "Application/json; charset=utf-8")
	str, err := result.String()
	fmt.Println("注销..", cli.ID, cli.App, cli.Uid, cli.Lease, str, err)
}

func (cli *Client) RunExample() {
	cli.Tls()
	mr.Seed(time.Now().UnixNano())
	n := mr.Intn(5000)
	time.Sleep(time.Duration(n+3000) * time.Millisecond)
	cli.POST() // 认证

	for i := 0; i < 800000000000000; i++ {
		mr.Seed(time.Now().UnixNano())
		n := mr.Intn(5000)
		time.Sleep(time.Duration(n+3000) * time.Millisecond)
		cli.PUT() // 心跳
	}
	cli.DELETE() // 退出
}

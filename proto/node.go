package proto

import (
	"context"
	"encoding/json"
	corec "github.com/offer365/example/grpc/core/client"
	cores "github.com/offer365/example/grpc/core/server"
	"github.com/offer365/example/tools"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"

	"github.com/zcalusic/sysinfo"
	"time"
)

const (
	Server_crt = `-----BEGIN CERTIFICATE-----
MIIDAzCCAesCCQDqO1aVDNi/IzANBgkqhkiG9w0BAQsFADBDMQswCQYDVQQGEwJH
QjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazETMBEGA1UEAwwKZ2l0
aHViLmNvbTAgFw0xOTEwMjAxMTAyMzhaGA8yMTE5MDkyNjExMDIzOFowQjELMAkG
A1UEBhMCR0IxDjAMBgNVBAcMBUNoaW5hMQ8wDQYDVQQKDAZzZXJ2ZXIxEjAQBgNV
BAMMCXNlcnZlci5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMRF
1mgbKYNO2X0iqX89Rrzc+xftqegQ+7V0n9Sa1HE07xQVcgu05faGB4B/29HPQ/gh
JMt1IxkXlISNuQwIDM5XVSahkH1OhQmtQnTKjYXFgboRFHMQUk26lKoIZ3o9AJ8s
QTPCLBw7a9StBpeWhBzEumDymP60hmGhTft4tbY85MrmObfTZ8KbQiHIy22jXNGV
N5ok61q4tlMV8HYK89q4WX7IcQusdK9NNwL1jZNQ4+WICEe2/zs8xY9r4REONKoM
HOME5aS+EvQSVwh5LyvNuPxa8io83EjokT3yRqZllvmXD/hVS/BCM927fgsiDfm0
ezuE5+AGiMR1N0agv1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAHxXVd/v7noVZ
LJ8IsLty3BjMX7ZjSvkyyrchxdQQIfCoMc/UGkDZ5TgvdPkE8eAfdSVwwrcpGf8V
C4ccB9flekd6HrO7Uo9mWrKcjyQn2MjwAZNDhcs5Sxrz8TusJEQk4iQYSq0oc4Nr
qGrR/2kXEirwXi/xQ0saVXalfhkK5W+rO/YWTc8K3leARQ6BDjGbF2BHRtj6HEZL
RnhJEbx+BvplXMlWQ5CBBYt/NQa/MKJDd2stT70Si8E1lIGIGaVQAy43uT7xy8XW
jSrruOAv1SVLovhSxjsMiu/jXwZsVAtaFAuT4ajiWQHzbNqUjVnt7dJIWJPCnL6h
lhaV0MUy7Q==
-----END CERTIFICATE-----

`
	Server_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAxEXWaBspg07ZfSKpfz1GvNz7F+2p6BD7tXSf1JrUcTTvFBVy
C7Tl9oYHgH/b0c9D+CEky3UjGReUhI25DAgMzldVJqGQfU6FCa1CdMqNhcWBuhEU
cxBSTbqUqghnej0AnyxBM8IsHDtr1K0Gl5aEHMS6YPKY/rSGYaFN+3i1tjzkyuY5
t9NnwptCIcjLbaNc0ZU3miTrWri2UxXwdgrz2rhZfshxC6x0r003AvWNk1Dj5YgI
R7b/OzzFj2vhEQ40qgwc4wTlpL4S9BJXCHkvK824/FryKjzcSOiRPfJGpmWW+ZcP
+FVL8EIz3bt+CyIN+bR7O4Tn4AaIxHU3RqC/VwIDAQABAoIBAAQbHeghoVWw4ZXf
ksIpqwAqc0pF24cSS+G45dsRvh38KIA4DqG2EBV/KksC4bta5aYcM2PaOHi+6Il5
WYSp6nKqmwpq2NX2PYw9RqWg0yMYRaV50/6wObiMja2c7WU+P3QU/ewyRK/2gkP5
tqiXKn5bkzaR/KdfaWxDbpkzJkIArLAELqEBuS0noxikrfypPanGnXk7IDhYo+rZ
cE0UHOhpkeo7gXeVc9tU/cjTRwBK7awKLIDWyknHGrL28nxMqKf06SzxG2oz6Hn3
twOtwAUS7tjophOZ6WCStgCOVFf0Ue6yJmja9xgWy/r2jJsH5/VV0xJZvmWGxr8T
IQh4oskCgYEA92Katy0Cvl1kS1/cf0ExMtOzXIwtCDu35axGl1FR3VMcoboPmH2h
HrRxSpcIgkRXz7wxsj3zttBXu8assjmwtCWzbDIE0YGYQ3v1CwDITihAyhevhW4b
UxN181RhMo1qHIcgULsVR5+P857FAHRSWWewh77ZK7x17fdQJshZujMCgYEAyxuT
R1CthfC7rbIX359tD9jb1XtG+XCgygZYv+6uoknmWMMmUqgDmQ3u8p4kuHudB6gm
/kZkxrluwJM5B8UKC1NRkejHP2ZO8ygpEGQp7t1H3BBFSfUVlu+YmfD5SjHhK9U5
2t+hfyuO8m0r+XdYk6lliEYufVlPMzJffT3rSk0CgYEAhs+jRGMw9ZBrUXAB9w8N
wob/XVW+TJhOlMiXB2r3U8cw+SktyonbvaHTgzRfHK4ltDz4UAvWvi83QEr6XX12
wBUze6ieW5Vl5pCsbryUa5MgC4Fw0yO3nEQkqN+4wBW0V6uDfrsU050ukzJYZPD+
113cI31rV5wyH+YANcJEs2UCgYEAmh0SY8qT4E4KGoJIGyadWqjyJcqk0CDl4GVw
cjJp0DrCzhdFvPI/yKMJ7I6Szmj9fhHZhJdlYGTT5MvROlQIiw9tlYlLpo+62EZg
4k8egmDlZdXyvWt6Nk0XPbfbcLDoapogjDOkFxq2HL054NDuJR0kLYMTQ4nAztgq
HJ4fKwECgYEAinsJM6lw9m3eyRRuPRFE4jNwg5KmZRjVuZ06+nPW/Sb7GXdN+5e6
62y87e63MRTm1r2C4g/esnqAOcS6iRHQtdTFrG8DU/j9F+uaB5TWZTroxqQ6h2F0
OjGZcdCMohluWRztbas01OZKSoDx1pEfP+H4kKFJ8LhWQXLU0lWibEw=
-----END RSA PRIVATE KEY-----

`
	Client_crt = `-----BEGIN CERTIFICATE-----
MIIDAzCCAesCCQDqO1aVDNi/JDANBgkqhkiG9w0BAQsFADBDMQswCQYDVQQGEwJH
QjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazETMBEGA1UEAwwKZ2l0
aHViLmNvbTAgFw0xOTEwMjAxMTAyMzhaGA8yMTE5MDkyNjExMDIzOFowQjELMAkG
A1UEBhMCR0IxDjAMBgNVBAcMBUNoaW5hMQ8wDQYDVQQKDAZjbGllbnQxEjAQBgNV
BAMMCWNsaWVudC5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAObU
dbPa2MMnp7X0P6TeUM9+gJgRgVdrOm05EPnf4p1xEmFq09bGupZpD+pVoU/yH/oi
wA4gwYtgk5ETtTfTbF8YUma6LYDye2m98zXiyVWpTs9pmxVRUTcpnjmyIS7mXSNE
hShN26OCTk8DtlL9STFnFWQY2Sb9PVjwDWTrXkHalQU3PFEmoQ/QPbTbBN2gydDn
WkK6LxgTaSA9xMw/j5upZh58aoLVwd8IevzKn/YnwQBEC0ptVQGl6B5EUKabhTWh
q6c4gDAhcqhdRFZa4UMcOZnzgwEuR7XzJlTL3AanBXJu5sUjDPTweOhENcdSHBQ/
sX6Cr9NFRm6bQqOrjmsCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAD4IIkxuNITUM
bHU2ebLEMPq8Udhcl9s3mBlaWf3ecDi4Yu+MBTy+ggcRhnq7zqaRVaRxdhtyVVIA
3hFwrWZK38jPGKrI9qZLoqQJu3RFq241jOjVol6zAkYhuqvO8n9AKhShjUFHPfA+
TN1BC8qb30lwYZnELaHdKFM16f7uska2lMY6N8uYqySNWFz/B77zIqUACRnvyGfS
gJ8QRDcGAjA0+SEMKtI0tB0qsL4c+de8uPaUjyO5uzWLXJap50gBi5Zx17YE8aMk
wntulWmvYwJSokZIOVi+3PDSc8Zh2ukhddA716NF2U8c7N++BZFBGcGLZK86Hh76
HayLryJQmA==
-----END CERTIFICATE-----

`
	Client_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA5tR1s9rYwyentfQ/pN5Qz36AmBGBV2s6bTkQ+d/inXESYWrT
1sa6lmkP6lWhT/If+iLADiDBi2CTkRO1N9NsXxhSZrotgPJ7ab3zNeLJValOz2mb
FVFRNymeObIhLuZdI0SFKE3bo4JOTwO2Uv1JMWcVZBjZJv09WPANZOteQdqVBTc8
USahD9A9tNsE3aDJ0OdaQrovGBNpID3EzD+Pm6lmHnxqgtXB3wh6/Mqf9ifBAEQL
Sm1VAaXoHkRQppuFNaGrpziAMCFyqF1EVlrhQxw5mfODAS5HtfMmVMvcBqcFcm7m
xSMM9PB46EQ1x1IcFD+xfoKv00VGbptCo6uOawIDAQABAoIBAQCTFY5qrGiy8fHL
33cudvrHPLR0MbNZINp5/oLytdaQvBwaNxgFI1yBuzCJAUdoyb/Wg44dcoHhbgiZ
yRUQHYhQkA7xpnCYWeqJ1p/DFl90Vg4B3CkVzFsT61EHMpoyaFewwViX9gSei8ma
T6M9/mdFM4pN3geA8JzGry/ZvqCxFID3Sz4/08zq9UjS54GiZgJb3lyGazdDk3Gf
h2NukbBRtvdh8iILjEM38czgqTBrDqXlFa5q0p9oq+UPn9twaVZcJ9t4IrcWIgaD
l9cYRE/agXj0cRO/IVOi/RB0e/NLiR0XqXSo4Rx7uGcSJys1yuPt96OGMIh2+c99
VGJbzBsBAoGBAP8qagFe5kJNrjweo9yNhs0H/TFx1mhCqQKPNFouDtttaCDcNvXx
d3B5KYKgWpTJPaZ1eGfPeA4OTLhKCLVG7EVQXUUsztyDS1JpuUJkm1texA30g0sw
UWhLfQfFEgWCaIQkbqZv1J5OYrc2xvPqjHfP+NG1CAte1w5QQ7FA541fAoGBAOeV
rO0yF30sDOUXlixfKN35j2FIgVB0DOT6nkpPyh1OYcdcshu3utGqmOiN7twqwyiL
m3Uucix/JbTb2m+HSAX9/s/SHHOoXeUp21wVSGYesknrBEZt3VifINzu/OFCjLk1
Plx4F0am0WrsDnAtQwgpCV29lgQjmFsXQZlUW051AoGBAJdvpbAgkUmCbsixapCn
0fv3JNZmeFgyT7n8IZbvxNOHkAgIifnXEArJbdBfuMKa2KLlDsuVfuvgormw/pAP
goP0mRZH7JFEvrwvkMqNiQJmMLcTiaRjDb13J8InvHVWmw7pzF2s+yPk44NW2CbE
6g7leAeFiDuvUrTk//e/zGzDAoGBAJx4TLaWubghIzVGkni4cuxHydB5JKYvQucT
Tg/3iR/z7ay9vLltkhRHp7i47UJkwieK7CZok0vtPJTOVvAz/z3NN3VDCWY7w/Uq
KsQ0vQ4Cf4Ph/ql3Ya6XFaUw9Dtes6YPi2r+2PsriyMrCzZP3pKM538msU1qn24s
cG4gyPBhAoGBAL+VTkIaLK07qChlT0Y2hwbmfLwAlOrPguJps7D0C0aBUDPXylOO
S5myV8jp+htbP6Mn5MEzZHhvoVSEe9GiCv9E5KMisJjPtQRRRKGNPAnTt9KJVQ1U
BCggzbZzimK/EFR72woV0071B93C4jO07jEmvkCb3gzmyWkgjREZQusj
-----END RSA PRIVATE KEY-----

`
	Ca_crt = `-----BEGIN CERTIFICATE-----
MIIDWzCCAkOgAwIBAgIJAKst9d2m1o1TMA0GCSqGSIb3DQEBCwUAMEMxCzAJBgNV
BAYTAkdCMQ4wDAYDVQQHDAVDaGluYTEPMA0GA1UECgwGZ29ib29rMRMwEQYDVQQD
DApnaXRodWIuY29tMCAXDTE5MTAyMDExMDIzOFoYDzIxMTkwOTI2MTEwMjM4WjBD
MQswCQYDVQQGEwJHQjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazET
MBEGA1UEAwwKZ2l0aHViLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBALh9e8OuCRP0zMYjbCqUk5b+d6J3tC9INL3P0VwcmWx5jCpUQLz0SGafnIL8
LworJfQkbDdOKNol9Zt4vzsxdV1k2VaZuAY0qWG5Kg+n1tCml46By9mQH3B6ngKe
cNDdBmRGYYDkuqI9g8UBgRYT4TbIQJ1Ns4wuKQR02/kCUfWypvE+8bEQEXTRKcHo
inILmcO7RvhWkfwWVbTpUv7M7K8wwIGKawDgl3DeW5g+tss0PD/iCdMo0DMRHykx
4KeTsrPYdxpxgf42LwG0aJ+/28GzYCQ4mYJaTADr5pp+vlUZWtYK8m7fFXbpGlrU
5aLTA5aEPdIuyTa2/DZXl4JBxTkCAwEAAaNQME4wHQYDVR0OBBYEFCikHb0Ms/7f
jci0C5Amwvf7cFmYMB8GA1UdIwQYMBaAFCikHb0Ms/7fjci0C5Amwvf7cFmYMAwG
A1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAG+jH3wHkTqx97/9voaftE/b
0tkbV+9b3SxPv5KoW0fm24x6UNrMPE9APt0J00Vlv20LNc/tOWquyKGDIhhe29/x
ehte/l7doGVW0Wg3xQtiIT5aJdMHNy+bSLogzV5D5sbHcPStKNj3M1wwhMj03YZ7
Nj5ua/c5aUU+MBMv0C/FNPnB+GSeRO2MxYHsZP2mBEJaLhPZ+K29kwGPCVWIESCH
IOS/jew/kfpPLavuvyPqoGAfc1xpe6QQXZUEGCtzTDU/rl/hQWMxCJg85E1S5Usx
gahmAgIzeyFCjb2txOo65VtLM0DfzzkIX2PrLz7CyiXP40m8uBMtCDG+IZS0arQ=
-----END CERTIFICATE-----
`
)

var (
	Self          *Node
	auth          *Authentication
	ClientConns   *CliConns
	StaterClients  *StaterClis
	_username     = "C205v406x68f5IM7"
	_password     = "c9bJ3v7FQ11681EP"
)


func init() {
	Self = NewNode(config.Cfg.Name, config.Cfg.LocalGRpcAddr())
	ClientConns   =new(CliConns)
	StaterClients  =new(StaterClis)
	auth = &Authentication{
		User:     _username,
		Password: _password,
	}
}

func NewNode(name, addr string) (n *Node) {
	n = new(Node)
	n.Hardware = new(Hardware)
	n.Hardware.Networks = make([]*Network, 0)
	n.Hardware.Mem = new(Mem)
	n.Hardware.Cpu = new(Cpu)
	n.Hardware.Bios = new(Bios)
	n.Hardware.Host = new(Host)
	n.Hardware.Board = new(Board)
	n.Hardware.Product = new(Product)
	n.Attrs = new(Attrs)
	n.Attrs.Addr = addr
	n.Attrs.Name = name
	n.Attrs.Start = time.Now().Unix()
	n.Hardware.hw()
	return
}

func (hd *Hardware) hw() {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	hd.Host.Machineid = si.Node.MachineID
	hd.Host.Architecture = si.OS.Architecture
	hd.Host.Hypervisor = si.Node.Hypervisor

	hd.Product.Name = si.Product.Name
	hd.Product.Serial = si.Product.Serial
	hd.Product.Vendor = si.Product.Vendor

	hd.Board.Name = si.Board.Name
	hd.Board.Vendor = si.Board.Vendor
	hd.Board.Serial = si.Board.Serial

	hd.Bios.Vendor = si.BIOS.Vendor

	hd.Cpu.Vendor = si.CPU.Vendor
	hd.Cpu.Threads = uint32(si.CPU.Threads)
	hd.Cpu.Model = si.CPU.Model
	hd.Cpu.Cache = uint32(si.CPU.Cache)
	hd.Cpu.Cores = uint32(si.CPU.Cores)
	hd.Cpu.Cpus = uint32(si.CPU.Cpus)
	hd.Cpu.Speed = uint32(si.CPU.Speed)

	hd.Mem.Speed = uint32(si.Memory.Speed)
	hd.Mem.Type = si.Memory.Type

	hd.Networks = make([]*Network, 0)
	for _, val := range si.Network {
		nw := new(Network)
		nw.Speed = uint32(val.Speed)
		nw.Macaddress = val.MACAddress
		nw.Driver = val.Driver
		hd.Networks = append(hd.Networks, nw)
	}

}

func (n *Node) md5() {
	if n.Attrs.Hwmd5 == "" {
		byt, err := json.Marshal(n.Hardware)
		if err != nil {
			return
		}
		n.Attrs.Hwmd5 = tools.Md5sum(byt, nil)
	}
}

func (n *Node) Status(ctx context.Context, args *Args) (*Node, error) {
	n.Hardware.hw()
	n.md5()
	n.Attrs.Now = time.Now().Unix()
	return n,nil
}

func GetAllNodes(ctx context.Context) (nodes map[string]*Node) {
	var lock sync.Mutex
	var wait sync.WaitGroup
	//var  value atomic.Value
	nodes = make(map[string]*Node, 0)
	//value.Store(nodes)
	peers:=config.Cfg.AllGRpcAddr()
	wait.Add(len(peers))
	for remoteName, remoteAddr := range peers {
		go func(remoteN string, remoteA string) {
			defer wait.Done()
			if  remoteN == Self.Attrs.Name || remoteA == Self.Attrs.Addr{
				Self.Attrs.Now = time.Now().Unix()
				// 重新获取硬件信息
				Self.Hardware.hw()
				Self.md5()
				return
			}
			cli, ok := StaterClients.Get(remoteN)
			if ok && cli != nil {
				n, err := cli.Status(ctx, &Args{Name: Self.Attrs.Name, Addr: Self.Attrs.Addr}, grpc.WaitForReady(true))
				if err != nil {
					if conn,ok:=ClientConns.Get(remoteN); conn!=nil && ok{
						conn.Close()
						ClientConns.Del(remoteN)
						StaterClients.Del(remoteN)
					}
					NodeGRpcClient(remoteN, remoteA)
					log.Sugar.Errorf("node rpc dial %s %s failed. error: %v", remoteN, remoteA, err)
					return
				}
				log.Sugar.Infof("node rpc dial %s %s success.", remoteN, remoteA)
				lock.Lock()
				nodes[remoteN] = n
				lock.Unlock()
			}
			return
		}(remoteName, remoteAddr)
	}
	wait.Wait()
	nodes[config.Cfg.Name] = Self
	//sort.Slice(nodes, func(i, j int) bool {
	//	return nodes[i].name < nodes[j].name
	//})
	return
}

func AllNodeGRpcClient(peers map[string]string) {
	for name, addr := range peers {
		if name != Self.Attrs.Name && addr != Self.Attrs.Addr {
			NodeGRpcClient(name, addr)
		}
	}
}

func NodeGRpcClient(name, addr string) {
	var (
		Con *grpc.ClientConn
		err error
	)

	Con, err = corec.NewRpcClient(
		corec.WithAddr(addr),
		corec.WithDialOption(grpc.WithPerRPCCredentials(auth)),
		corec.WithServerName("server.io"),
		corec.WithCert([]byte(Client_crt)),
		corec.WithKey([]byte(Client_key)),
		corec.WithCa([]byte(Ca_crt)),
	)
	if err != nil {
		log.Sugar.Error(err)
		return
	}

	StaterClients.Put(name,NewStaterClient(Con))
	ClientConns.Put(name,Con)
	return
}

func NodeGRpcServer() (*grpc.Server, error) {
	// Token认证
	auth := func(ctx context.Context) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing credentials")
		}

		var user string
		var pwd string

		if val, ok := md["user"]; ok {
			user = val[0]
		}
		if val, ok := md["password"]; ok {
			pwd = val[0]
		}

		if user != _username || pwd != _password {
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return nil
	}

	// 一元拦截器
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}

	// 实例化grpc Server
	return cores.NewRpcServer(
		cores.WithServerOption(grpc.UnaryInterceptor(interceptor)),
		cores.WithCert([]byte(Server_crt)),
		cores.WithKey([]byte(Server_key)),
		cores.WithCa([]byte(Ca_crt)),
	)
}

// Authentication 自定义认证
// 要实现对每个gRPC方法进行认证，需要实现grpc.PerRPCCredentials接口
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


type CliConns struct {
	sync.Map
}

func (cc *CliConns)Get(key string) (*grpc.ClientConn,bool) {
	val,exist:=cc.Load(key)
	if exist {
		cli,ok:= val.(*grpc.ClientConn)
		return cli,ok
	}
	return nil,exist
}

func (cc *CliConns)Put(key string,val *grpc.ClientConn){
	cc.Store(key,val)
}

func (cc *CliConns)Del(key string)  {
	cc.Delete(key)
}

type StaterClis struct {
	sync.Map
}

func (sc *StaterClis)Get(key string) (StaterClient,bool) {
	val,exist:=sc.Load(key)
	if exist {
		cli,ok:= val.(StaterClient)
		return cli,ok
	}
	return nil,exist
}

func (sc *StaterClis)Put(key string,val StaterClient){
	sc.Store(key,val)
}

func (sc *StaterClis)Del(key string)  {
	sc.Delete(key)
}

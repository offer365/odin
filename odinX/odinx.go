package odinX

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/offer365/odin/log"
)

const (
	logo = `
	             _   _        
	            | | (_)       
	  ___     __| |  _   _ __  
	 / _ \   / _' | | | | '_ \
	| (_) | | (_| | | | | | | |
	 \___/   \__,_| |_| |_| |_|
	`
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(logo)
}

var (
	Cfg *Config
)

// 加解密方法
type CryptFunc func(src []byte) ([]byte, error)

// hash 方法
type HashFunc func(byt []byte) string

type Config struct {
	// etcd embed config
	EmbedCtx          context.Context
	EmbedName         string
	EmbedDir          string
	EmbedClientAddr   string
	EmbedPeerAddr     string
	EmbedClusterToken string
	EmbedClusterState string
	EmbedCluster      map[string]string
	EmbedAuthPwd      string

	// etcd client config
	EtcdCliCtx     context.Context
	EtcdCliAddr    string
	EtcdCliUser    string
	EtcdCliTimeout time.Duration

	// store key
	StoreLicenseKey            string
	StoreClearLicenseKey       string
	StoreClientConfigKeyPrefix string
	StoreClientKeyPrefix       string
	StoreTokenKey              string
	StoreSerialNumKey          string

	// gRpc config
	GRpcServerCrt  string
	GRpcServerKey  string
	GRpcClientCrt  string
	GRpcClientKey  string
	GRpcCaCrt      string
	GRpcUser       string
	GRpcPwd        string
	GRpcServerName string
	GRpcAllNode    map[string]string
	GRpcListen     string
	RestfulPwd     string

	// node config
	NodeName     string
	NodeAddr     string
	NodeHardware HardWare

	// encrypt decrypt func
	// odin & edda
	LicenseEncrypt CryptFunc // license 加解密
	LicenseDecrypt CryptFunc
	SerialEncrypt  CryptFunc // 序列号 加解密
	SerialDecrypt  CryptFunc
	UntiedEncrypt  CryptFunc // 解绑码 加解密
	UntiedDecrypt  CryptFunc

	// odin & client app
	VerifyDecrypt CryptFunc // token 密文解密
	CipherEncrypt CryptFunc // uuid 加密
	AuthEncrypt   CryptFunc // auth 数据加密
	UuidHash      HashFunc
	TokenHash     HashFunc
}

func NewConfig() *Config {
	return &Config{
		EmbedCtx:          context.TODO(),
		EmbedName:         "",
		EmbedDir:          "disk/default",
		EmbedClientAddr:   "127.0.0.1:12379",
		EmbedPeerAddr:     "127.0.0.1:12380",
		EmbedClusterToken: "",
		EmbedClusterState: "new",
		EmbedCluster:      nil,
		EmbedAuthPwd:      "",

		EtcdCliCtx:  context.TODO(),
		EtcdCliAddr: "127.0.0.1:12379",
		EtcdCliUser: "root",

		EtcdCliTimeout: 3 * time.Second,

		StoreLicenseKey:            "/Mzot24SEaI91buzPv8C74302O905dc13gP7bXo69QSP1ot9Wt7BiET6YN7Zah72S/p5353Ls032rcUbJD8759w059Q70fyjOLGl0dblSDEI2RrTC3t5rHA5hJ540d9504",
		StoreClearLicenseKey:       "/377B3Rs9A24J5wMX1WCRt8ANtqZh85K9xQ28H901o78jsze0b806P01Y2t1MW7Sf/vPQoG0xW18i93TDO6ec1y54Qz8lbb2tz2G5sG5aaVXKL31Ji927q8hT1v0eZH11V",
		StoreClientConfigKeyPrefix: "/ABw6T376rymCP8CC8c6Z0a5010xCRgk9AQJ8Y8H29e2mt42ng62fs915X8o9POSH/a1JF81wHq7BN57l3dS23oi5KkXz1JE30R328ys22Z1WoO2PiKZ9Mg00d23G6s18f/",
		StoreClientKeyPrefix:       "/ogyI0o7jJm0Pb07rAk8v5pI116lWbfJko03c33dYnKxu03i5x2N9j7XS3B6CJ86g/TLx1u38pWD89GdVmymUDQElB9l36dn65t37X1o747RV9c1PTY525216LN8UUN1sJ/",
		StoreTokenKey:              "/bZdsG90ST0m59evU7obGb2dqST11gVq4GRdjGD0HkBeS2Qh2v7FsOVXIOMDH9nRb/k3KIz7fWos1BV7SW5m88Vh8MYLaNFZ8lwKr09X8V4ewUJ87t55AdC5C0Gq8cRZeC/",
		StoreSerialNumKey:          "/L30GZduFHLqXZSbrvDWi90Ik87UG1Vsc5FDQZByvL9C3Ad9mD9DWEb17NLZM4dV3/Y4iohWvB55Ef1Sg6b1uxPtvM3rg71182B6wEW1lBnNCwPaKKmq74DBp65c998J25",

		GRpcServerCrt:  "",
		GRpcServerKey:  "",
		GRpcClientCrt:  "",
		GRpcClientKey:  "",
		GRpcCaCrt:      "",
		GRpcUser:       "",
		GRpcPwd:        "",
		GRpcServerName: "",
		GRpcAllNode:    nil,
		GRpcListen:     "0.0.0.0:9527",
		RestfulPwd:     "",

		NodeName:     "",
		NodeAddr:     "",
		NodeHardware: nil,

		LicenseEncrypt: nil,
		LicenseDecrypt: nil,
		SerialEncrypt:  nil,
		SerialDecrypt:  nil,
		UntiedEncrypt:  nil,
		UntiedDecrypt:  nil,
		VerifyDecrypt:  nil,
		CipherEncrypt:  nil,
		AuthEncrypt:    nil,
		UuidHash:       nil,
		TokenHash:      nil,
	}
}

func (cfg Config) CheckValue() (err error) {
	v := reflect.ValueOf(cfg)
	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch vt.Field(i).Type.Name() {
		case "string":
			if f.String() == "" {
				return errors.New(vt.Field(i).Name + " field cannot be empty")
			}
		case "Duration":
			if f.IsZero() {
				return errors.New(vt.Field(i).Name + " field cannot be zero")
			}
		default:
			if f.IsNil() {
				return errors.New(vt.Field(i).Name + " field cannot be nil")
			}
		}
	}
	return
}

func Start(cfg *Config)  {
	if err := cfg.CheckValue(); err != nil {
		log.Sugar.Fatal(err)
		return
	}
	Cfg = cfg
	Main()
	return
}

type HardWare interface {
	GetSysInfo()
	HostInfo() (machineID, architecture, hypervisor string)
	ProductInfo() (name, serial, vendor string)
	BoardInfo() (name, serial, vendor string)
	BiosInfo() (vendor string)
	CpuInfo() (vendor, model string, threads, cache, cores, cpus, speed uint32)
	MemInfo() (speed uint32, tp string)
	NetworksInfo() []*NetDriver
}

type NetDriver struct {
	Driver     string
	Macaddress string
	Speed      uint32
}

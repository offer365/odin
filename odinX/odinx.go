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

	// web config
	WebPwd    string
	WebListen string

	// node config
	NodeName     string
	NodeAddr     string
	NodeHardware HardWare

	// encrypt decrypt func
	// odin & edda
	LicenseEncrypt CryptFunc // license 加解密
	LicenseDecrypt CryptFunc
	SerialEncrypt  CryptFunc // 序列号 加解密
	// SerialDecrypt  CryptFunc
	// UntiedEncrypt  CryptFunc // 解绑码 加解密
	UntiedDecrypt CryptFunc
	// ClearEncrypt   CryptFunc // 注销码 加解密
	// ClearDecrypt   CryptFunc
	TokenHash HashFunc

	// odin & client app
	VerifyDecrypt CryptFunc // token 密文解密
	CipherEncrypt CryptFunc // uuid 加密
	AuthEncrypt   CryptFunc // auth 数据加密
	UuidHash      HashFunc
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

		StoreLicenseKey:            "",
		StoreClearLicenseKey:       "",
		StoreClientConfigKeyPrefix: "",
		StoreClientKeyPrefix:       "",
		StoreTokenKey:              "",
		StoreSerialNumKey:          "",

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
		WebPwd:         "",

		NodeName:     "",
		NodeAddr:     "",
		NodeHardware: nil,

		LicenseEncrypt: nil,
		LicenseDecrypt: nil,
		SerialEncrypt:  nil,
		// SerialDecrypt:  nil,
		// UntiedEncrypt:  nil,
		UntiedDecrypt: nil,
		TokenHash:     nil,

		VerifyDecrypt: nil,
		CipherEncrypt: nil,
		AuthEncrypt:   nil,
		UuidHash:      nil,
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

func Start(cfg *Config) {
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
	BiosInfo() (vendor, version string)
	CpuInfo() (vendor, model string, threads, cache, cores, cpus, speed uint32)
	MemInfo() (speed uint32, tp string)
	NetworksInfo() []*NetDriver
}

type NetDriver struct {
	Driver     string
	Macaddress string
	Speed      uint32
}

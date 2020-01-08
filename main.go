package main

import (
	"context"
	"time"

	"github.com/offer365/example/endecrypt/endeaes"
	"github.com/offer365/example/endecrypt/endeaesrsa"
	"github.com/offer365/example/endecrypt/endeaesrsaecc"
	"github.com/offer365/example/endecrypt/endersa"
	"github.com/zcalusic/sysinfo"

	"github.com/offer365/odin/config"
	"github.com/offer365/odin/odinX"
	"github.com/offer365/odin/utils"
)

var hw = &hardware{}

func main() {
	cfg := &odinX.Config{
		EmbedCtx:          context.TODO(),
		EmbedName:         config.Cfg.Name,
		EmbedDir:          config.Cfg.Dir,
		EmbedClientAddr:   config.Cfg.LocalClientAddr(),
		EmbedPeerAddr:     config.Cfg.LocalPeerAddr(),
		EmbedClusterToken: clusterToken,
		EmbedClusterState: config.Cfg.State,
		EmbedCluster:      config.Cfg.AllPeerAddr(),
		EmbedAuthUser:     "root",
		EmbedAuthPwd:      embedAuthPwd,

		EtcdCliCtx:     context.TODO(),
		EtcdCliAddr:    config.Cfg.LocalClientAddr(),
		EtcdCliTimeout: 5 * time.Second,

		StoreLicenseKey:            storeLicenseKey,
		StoreClearLicenseKey:       storeClearLicenseKey,
		StoreClientConfigKeyPrefix: storeClientConfigKeyPrefix,
		StoreClientKeyPrefix:       storeClientKeyPrefix,
		StoreTokenKey:              storeTokenKey,
		StoreSerialNumKey:          storeSerialNumKey,

		GRpcServerCrt:  server_crt,
		GRpcServerKey:  server_key,
		GRpcClientCrt:  client_crt,
		GRpcClientKey:  client_key,
		GRpcCaCrt:      ca_crt,
		GRpcUser:       grpcUser,
		GRpcPwd:        grpcPwd,
		GRpcServerName: server_name,
		GRpcAllNode:    config.Cfg.AllGRpcAddr(),
		GRpcListen:     config.Cfg.LocalGRpcAddr(),


		// web config
		WebPwd:    config.Cfg.Pwd,
		WebListen: config.Cfg.Web,

		NodeName:     config.Cfg.LocalName(),
		NodeAddr:     config.Cfg.LocalGRpcAddr(),
		NodeHardware: hw,

		// odin & edda
		LicenseEncrypt: licenseEncrypt1,
		LicenseDecrypt: licenseDecrypt1,
		SerialEncrypt:  serialEncrypt1,
		// SerialDecrypt:              PubDecryptRsa2048Aes256,
		// UntiedEncrypt:              PriEncryptRsa2048Aes256,
		UntiedDecrypt: untiedDecrypt1,
		TokenHash:     HashFunc1,

		// odin & app
		VerifyDecrypt: verifyDecrypt1,
		CipherEncrypt: cipherEncrypt1,
		AuthEncrypt:   authEncrypt1,
		UuidHash:      HashFunc2,
	}
	odinX.Start(cfg)
}

// odin & edda

// Pub Encrypt Rsa2048 Aes256
func licenseEncrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub1), []byte(_aes256key1))
}

// Pri Decrypt Rsa2048 Aes256
func licenseDecrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri1), []byte(_aes256key1))
}

// Pub Encrypt Ecc256 Rsa204 8Aes256
func licenseEncrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub1), []byte(_rsa2048pub1), []byte(_aes256key1))
}

// Pri Decrypt Ecc25 6Rsa2048 Aes256
func licenseDecrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri1), []byte(_rsa2048pri1), []byte(_aes256key1))
}

// Pub Encrypt Rsa2048 Aes256
func serialEncrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub2), []byte(_aes256key2))
}

// Pri Decrypt Rsa2048 Aes256
func serialDecrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri2), []byte(_aes256key2))
}

// Pub Encrypt Ecc256 Rsa204 8Aes256
func serialEncrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub2), []byte(_rsa2048pub2), []byte(_aes256key2))
}

// Pri Decrypt Ecc25 6Rsa2048 Aes256
func serialDecrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri2), []byte(_rsa2048pri2), []byte(_aes256key2))
}

// Pub Encrypt Rsa2048 Aes256
func untiedEncrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub3), []byte(_aes256key3))
}

// Pri Decrypt Rsa2048 Aes256
func untiedDecrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri3), []byte(_aes256key3))
}

// Pub Encrypt Ecc256 Rsa204 8Aes256
func untiedEncrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub3), []byte(_rsa2048pub3), []byte(_aes256key3))
}

// Pri Decrypt Ecc25 6Rsa2048 Aes256
func untiedDecrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri3), []byte(_rsa2048pri3), []byte(_aes256key3))
}

func PriDecryptRsa2048(src []byte) ([]byte, error) {
	return endersa.PriDecrypt(src, []byte(_rsa4096pri1))
}

func Aes256key1(src []byte) ([]byte, error) {
	return endeaes.AesCbcEncrypt(src, []byte(_aes256key4))
}

func Aes256key2(src []byte) ([]byte, error) {
	return endeaes.AesCbcEncrypt(src, []byte(_aes256key4))
}

func HashFunc1(src []byte) string {
	return utils.Sha256Hex(src, []byte(storeHashSalt1))
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

type hardware struct {
	// linux
	// TODO: 当 sysinfo.SysInfo 为这个实例的全局变量时，在多次获取网卡与磁盘信息时，会导致返回结果不断堆叠。
	//  解决方法是，修改这个包在 getNetworkInfo for循环前增加 si.Network=make([]NetworkDevice,0) 磁盘同理。
<<<<<<< HEAD
	sysinfo.SysInfo  // "github.com/zcalusic/sysinfo"
=======
	// sysinfo.SysInfo  // "github.com/zcalusic/sysinfo"
>>>>>>> 2ae8a808efc2075d1a1a04ef06400585e3798b30
	// windows
	//winsysinfo.SysInfo // "github.com/offer365/example/winsysinfo"
}

func (h *hardware) HostInfo() (machineID, architecture, hypervisor string) {
	return h.Node.MachineID, h.OS.Architecture, h.Node.Hypervisor
}

func (h *hardware) ProductInfo() (name, serial, vendor string) {
	return h.Product.Name, h.Product.Serial, h.Product.Vendor
}

func (h *hardware) BoardInfo() (name, serial, vendor string) {
	return h.Board.Name, h.Board.Serial, h.Board.Vendor
}

func (h *hardware) BiosInfo() (vendor, version,date string) {
	return h.BIOS.Vendor, h.BIOS.Version,h.BIOS.Date
}

func (h *hardware) CpuInfo() (vendor, model string, threads, cache, cores, cpus, speed uint32) {
	return h.CPU.Vendor, h.CPU.Model, uint32(h.CPU.Threads),
		uint32(h.CPU.Cache), uint32(h.CPU.Cores),
		uint32(h.CPU.Cpus), uint32(h.CPU.Speed)
}

func (h *hardware) MemInfo() (speed uint32, tp string) {
	return uint32(h.Memory.Speed), h.Memory.Type
}

func (h *hardware) NetworksInfo() (drivers []*odinX.NetDriver) {
	for _, val := range h.Network {
		nw := new(odinX.NetDriver)
		nw.Speed = uint32(val.Speed)
		nw.Macaddress = val.MACAddress
		nw.Driver = val.Driver
		drivers = append(drivers, nw)
	}
	return
}

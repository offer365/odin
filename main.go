package main

import (
	"context"
	"time"

	"github.com/offer365/example/endecrypt/endeaes"
	"github.com/offer365/example/endecrypt/endeaesrsa"
	"github.com/offer365/example/endecrypt/endeaesrsaecc"
	"github.com/offer365/example/endecrypt/endersa"
	"github.com/offer365/example/winsysinfo"
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
		EmbedAuthPwd:      embedAuthPwd,

		EtcdCliCtx:                 context.TODO(),
		EtcdCliAddr:                "127.0.0.1:" + config.Cfg.LocalClientPort(),
		EtcdCliUser:                "root",
		EtcdCliTimeout:             3 * time.Second,
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
		LicenseEncrypt: PubEncryptRsa2048Aes256,
		LicenseDecrypt: PriDecryptRsa2048Aes256,
		SerialEncrypt:  PubEncryptRsa2048Aes256,
		// SerialDecrypt:              PubDecryptRsa2048Aes256,
		// UntiedEncrypt:              PriEncryptRsa2048Aes256,
		UntiedDecrypt: PubDecryptRsa2048Aes256,
		TokenHash:     HashFunc,

		// odin & app
		VerifyDecrypt: PriDecryptRsa2048,
		CipherEncrypt: Aes256key1,
		AuthEncrypt:   Aes256key2,
		UuidHash:      HashFunc,
	}
	odinX.Start(cfg)
}

func PubEncryptRsa2048Aes256(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub1), []byte(_aes256key1))
}

func PriDecryptRsa2048Aes256(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri1), []byte(_aes256key1))
}

// Ecc256 + Rsa2048 + Aes256
func PubEncryptEcc256Rsa2048Aes256(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub2), []byte(_rsa2048pub2), []byte(_aes256key2))
}

// Ecc256 + Rsa2048 + Aes256
func PriDecryptEcc256Rsa2048Aes256(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri2), []byte(_rsa2048pri2), []byte(_aes256key2))
}

func PubDecryptRsa2048Aes256(src []byte) ([]byte, error) {
	return endeaesrsa.PubDecrypt(src, []byte(_rsa2048pub3), []byte(_aes256key3))
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

func HashFunc(src []byte) string {
	return utils.Sha256Hex(src, []byte(storeHashSalt))
}

type hardware struct {
	// linux
	// sysinfo.SysInfo  // "github.com/zcalusic/sysinfo"
	// windows
	winsysinfo.SysInfo // "github.com/offer365/example/winsysinfo"
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

func (h *hardware) BiosInfo() (vendor string) {
	return h.BIOS.Vendor
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

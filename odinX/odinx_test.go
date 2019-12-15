package odinX

import (
	"fmt"
	"testing"
)

func TestCheckValue(t *testing.T) {
	cfg := NewConfig()
	cfg.EmbedName = "2"
	cfg.EmbedClusterToken = "2"
	cfg.EmbedCluster = map[string]string{}
	cfg.EmbedAuthPwd = "2"
	cfg.EtcdCliTimeout = 10
	cfg.GRpcServerCrt = "32"
	cfg.GRpcPwd = "3"
	cfg.GRpcAllNode = map[string]string{}
	cfg.GRpcCaCrt = "2"
	cfg.GRpcServerName = "3"
	cfg.GRpcServerKey = "s"
	cfg.GRpcCaCrt = "33"
	cfg.GRpcClientCrt = "gf"
	cfg.GRpcClientKey = "dd"
	cfg.GRpcUser = "ssf"
	cfg.NodeAddr = "sdf"
	cfg.NodeName = "fd"
	err := cfg.CheckValue()
	fmt.Println(err)
}

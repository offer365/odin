package config

import (
	"flag"
	"github.com/offer365/odin/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

var (
	Cfg *Config
	cfp string
)

func args() {
	flag.StringVar(&cfp, "f", "odin.yaml", "Config file path.")
	flag.Parse()
}

func init() {
	Cfg = new(Config)
	args()
	Cfg.LoadYaml(cfp)
}

type Config struct {
	Name  string          `yaml:"name" json:"name"`
	Peers map[string]Node `yaml:"peers" json:"peers"`
	Dir   string          `yaml:"dir" json:"dir"`
	State string          `yaml:"state" json:"state"`
	Pwd   string          `yaml:"pwd" json:"pwd"`
}

type Node struct {
	Peer   string `yaml:"peer" json:"peer"`
	Client string `yaml:"client" json:"client"`
	GRpc   string `json:"grpc" json:"grpc"`
}

func (cfg *Config) LoadYaml(filename string) {
	var (
		content []byte
		err     error
		//name    string
	)
	cfg.Peers = make(map[string]Node)
	//读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Sugar.Fatal("failed to read configuration file. error: ", err)
	}
	if err = yaml.Unmarshal(content, cfg); err != nil {
		log.Sugar.Fatal("failed to unmarshal configuration file. error: ", err)
	}
	return
}

func (cfg *Config) AllPeerAddr() (pps map[string]string) {
	pps = make(map[string]string)
	for name, node := range cfg.Peers {
		pps[name] = node.Peer
	}
	return
}

func (cfg *Config) AllGRpcAddr() (pps map[string]string) {
	pps = make(map[string]string)
	for name, node := range cfg.Peers {
		pps[name] = node.GRpc
	}
	return
}

func (cfg *Config) AllClientAddr() (pps map[string]string) {
	pps = make(map[string]string)
	for name, node := range cfg.Peers {
		pps[name] = node.Client
	}
	return
}

func (cfg *Config) LocalName() string {
	return cfg.Name
}

func (cfg *Config) LocalIp() string {
	return cfg.GetIp(cfg.Name)
}

func (cfg *Config) LocalGRpcAddr() string {
	return cfg.GetGRpcAddr(cfg.Name)
}

func (cfg *Config) LocalPeerAddr() string {
	return cfg.GetPeerAddr(cfg.Name)
}

func (cfg *Config) LocalClientAddr() string {
	return cfg.GetClientAddr(cfg.Name)
}

func (cfg *Config) GetIp(name string) string {
	node, ok := cfg.Peers[name]
	if !ok {
		return ""
	}
	lis := strings.Split(node.Peer, ":")
	if len(lis) >= 1 {
		return lis[0]
	}
	return ""
}

func (cfg *Config) GetGRpcAddr(name string) string {
	node, ok := cfg.Peers[name]
	if !ok {
		return ""
	}
	return node.GRpc
}

func (cfg *Config) GetPeerAddr(name string) string {
	node, ok := cfg.Peers[name]
	if !ok {
		return ""
	}
	return node.Peer
}

func (cfg *Config) GetClientAddr(name string) string {
	node, ok := cfg.Peers[name]
	if !ok {
		return ""
	}
	return node.Client
}

func (cfg *Config) LocalClientPort() string {
	return cfg.GetClientPort(cfg.Name)
}

func (cfg *Config) LocalGRpcPort() string {
	return cfg.GetGRpcPort(cfg.Name)
}

func (cfg *Config) LocalPeerPort() string {
	return cfg.GetPeerPort(cfg.Name)
}

func (cfg *Config) GetGRpcPort(name string) string {
	lis := strings.Split(cfg.GetGRpcAddr(name), ":")
	if len(lis) > 1 {
		return lis[1]
	}
	return ""
}

func (cfg *Config) GetPeerPort(name string) string {
	lis := strings.Split(cfg.GetPeerAddr(name), ":")
	if len(lis) > 1 {
		return lis[1]
	}
	return ""
}

func (cfg *Config) GetClientPort(name string) string {
	lis := strings.Split(cfg.GetClientAddr(name), ":")
	if len(lis) > 1 {
		return lis[1]
	}
	return ""
}

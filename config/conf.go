package config

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/offer365/odin/log"
	"gopkg.in/yaml.v2"
)

var (
	Cfg *config
	cfp string
)

func args() {
	flag.StringVar(&cfp, "f", "odin.yaml", "config file path.")
	flag.Parse()
}

func init() {
	Cfg = new(config)
	args()
	Cfg.LoadYaml(cfp)
}

type config struct {
	Name  string          `yaml:"name" json:"name"`
	Peers map[string]Node `yaml:"peers" json:"peers"`
	Dir   string          `yaml:"dir" json:"dir"`
	State string          `yaml:"state" json:"state"`
	Web   string          `yaml:"web" json:"web"`
	Pwd   string          `yaml:"pwd" json:"pwd"`
}

type Node struct {
	Peer   string `yaml:"peer" json:"peer"`
	Client string `yaml:"client" json:"client"`
	GRpc   string `yaml:"grpc" json:"grpc"`
}

func (c *config) LoadYaml(filename string) {
	var (
		content []byte
		err     error
		// name    string
	)
	c.Peers = make(map[string]Node)
	// 读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Sugar.Fatal("failed to read configuration file. error: ", err)
	}
	if err = yaml.Unmarshal(content, c); err != nil {
		log.Sugar.Fatal("failed to unmarshal configuration file. error: ", err)
	}
	return
}

func (c *config) AllPeerAddr() (pps map[string]string) {
	pps = make(map[string]string)
	for name, node := range c.Peers {
		pps[name] = node.Peer
	}
	return
}

func (c *config) AllGRpcAddr() (pps map[string]string) {
	pps = make(map[string]string)
	for name, node := range c.Peers {
		pps[name] = node.GRpc
	}
	return
}

func (c *config) AllClientAddr() (pps map[string]string) {
	pps = make(map[string]string)
	for name, node := range c.Peers {
		pps[name] = node.Client
	}
	return
}

func (c *config) LocalName() string {
	return c.Name
}

func (c *config) LocalIp() string {
	return c.GetIp(c.Name)
}

func (c *config) LocalGRpcAddr() string {
	return c.GetGRpcAddr(c.Name)
}

func (c *config) LocalPeerAddr() string {
	return c.GetPeerAddr(c.Name)
}

func (c *config) LocalClientAddr() string {
	return c.GetClientAddr(c.Name)
}

func (c *config) GetIp(name string) string {
	node, ok := c.Peers[name]
	if !ok {
		return ""
	}
	lis := strings.Split(node.Peer, ":")
	if len(lis) >= 1 {
		return lis[0]
	}
	return ""
}

func (c *config) GetGRpcAddr(name string) string {
	node, ok := c.Peers[name]
	if !ok {
		return ""
	}
	return node.GRpc
}

func (c *config) GetPeerAddr(name string) string {
	node, ok := c.Peers[name]
	if !ok {
		return ""
	}
	return node.Peer
}

func (c *config) GetClientAddr(name string) string {
	node, ok := c.Peers[name]
	if !ok {
		return ""
	}
	return node.Client
}

func (c *config) LocalClientPort() string {
	return c.GetClientPort(c.Name)
}

func (c *config) LocalGRpcPort() string {
	return c.GetGRpcPort(c.Name)
}

func (c *config) LocalPeerPort() string {
	return c.GetPeerPort(c.Name)
}

func (c *config) GetGRpcPort(name string) string {
	lis := strings.Split(c.GetGRpcAddr(name), ":")
	if len(lis) > 1 {
		return lis[1]
	}
	return ""
}

func (c *config) GetPeerPort(name string) string {
	lis := strings.Split(c.GetPeerAddr(name), ":")
	if len(lis) > 1 {
		return lis[1]
	}
	return ""
}

func (c *config) GetClientPort(name string) string {
	lis := strings.Split(c.GetClientAddr(name), ":")
	if len(lis) > 1 {
		return lis[1]
	}
	return ""
}

package config

import (
	"encoding/json"
	"github.com/offer365/odin/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

var (
	ConfFilePath string
	Cfg          = new(Config)
)

type Config struct {
	Name  string
	Group string
	Nodes
	Ports
	Dir   string `json:"dir"`
	State string `json:"state"`
	Pwd   string `json:"pwd"`
}

type Nodes struct {
	Addr  string   `json:"addr"`
	Peers []string `json:"peers"`
}

type Ports struct {
	Peer    string `json:"peer"`
	Client  string `json:"client"`
	Web     string `json:"web"`
	Rpc     string `json:"rpc"`
	Metrics string `json:"metrics"`
}

func (cfg *Config) LoadJson(filename string) {
	var (
		content []byte
		err     error
	)
	//读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Sugar.Fatal("failed to read configuration file. error: ", err)
	}
	// json反序列化
	if err = json.Unmarshal(content, cfg); err != nil {
		log.Sugar.Fatal("failed to unmarshal configuration file. error: ", err)
	}

	for id, ip := range cfg.Peers {
		if ip == cfg.Addr {
			cfg.Name = cfg.Group + strconv.Itoa(id)
		}
	}
	return
}

func (cfg *Config) LoadYaml(filename string) {
	var (
		content []byte
		err     error
	)
	//读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Sugar.Fatal("failed to read configuration file. error: ", err)
	}
	// json反序列化
	if err = yaml.Unmarshal(content, cfg); err != nil {
		log.Sugar.Fatal("failed to unmarshal configuration file. error: ", err)
	}

	for id, ip := range cfg.Peers {
		if ip == cfg.Addr {
			cfg.Name = cfg.Group + strconv.Itoa(id)
		}
	}
	return
}

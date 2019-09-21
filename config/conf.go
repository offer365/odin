package config

import (
	"../log"
	"encoding/json"
	"flag"
	"io/ioutil"
	"strconv"
)

var (
	ConfFilePath string
	Cfg          = new(Config)
)

type Config struct {
	Name   string
	Addr   string   `json:"addr"`
	Peer   string   `json:"peer"`
	Client string   `json:"client"`
	Web    string   `json:"web"`
	Rpc    string   `json:"rpc"`
	//Ntp    bool     `json:"ntp"`
	Dir    string   `json:"dir"`
	Peers  []string `json:"peers"`
	State  string   `json:"state"`
	Pwd    string   `json:"pwd"`
}

func (cfg *Config) LoadConfig(filename string) {
	var (
		content []byte
		err     error
	)
	//读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		goto ERR
	}
	// json反序列化
	if err = json.Unmarshal(content, cfg); err != nil {
		goto ERR
	}

	for id, ip := range cfg.Peers {
		if ip == cfg.Addr {
			cfg.Name = "odin" + strconv.Itoa(id)
		}
	}
	return
ERR:
	cfg.Addr = "127.0.0.1"
	cfg.Peers = []string{"127.0.0.1"}
	cfg.Peer = "12380"
	cfg.Client = "12379"
	log.Sugar.Error("Failed to load configuration file. error: ", err.Error())
	return
}

func args() {
	flag.StringVar(&ConfFilePath, "f", "odin.json", "Config file path.")
	flag.Parse()
}
func init() {
	args()
	Cfg.LoadConfig(ConfFilePath)
}

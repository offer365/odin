package odinmain

import (
	"flag"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/offer365/odin/asset"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/logic"
	"github.com/offer365/odin/model"
	"github.com/offer365/odin/node"
	"strings"

	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	Username = "root"
	Password = "613f#8d164df4ACPF49@93a510df49!66f98b*d6"
	logo     = `
	             _   _        
	            | | (_)       
	  ___     __| |  _   _ __  
	 / _ \   / _' | | | | '_ \
	| (_) | | (_| | | | | | | |
	 \___/   \__,_| |_| |_| |_|
	`
)

var (
	AssetPath string
	User      = "admin"
	debug     bool
	cfp       string
)

func args() {
	flag.StringVar(&cfp, "f", "odin.yaml", "Config file path.")
	flag.Parse()
}

// 释放静态资源
func RestoreAsset() {
	// 解压 静态文件的位置
	if runtime.GOOS == "linux" {
		AssetPath = "/usr/share/.asset/.temp/"
	} else {
		AssetPath = "./"
	}
	// go get -u github.com/jteeuwen/go-bindata/...
	// 重新生成静态资源在项目的根目录下 go-bindata -o=asset/asset.go -pkg=asset html/... static/...
	dirs := []string{"html", "static"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets(AssetPath, dir); err != nil {
			log.Sugar.Error("restore assets failed. error: ", err)
			_ = os.RemoveAll(filepath.Join(AssetPath, dir))
			continue
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(logo)
	args()
	RestoreAsset()
	switch {
	case strings.HasSuffix(cfp, ".json"):
		config.Cfg.LoadJson(cfp)
	case strings.HasSuffix(cfp, ".yaml"):
		config.Cfg.LoadYaml(cfp)
	default:
		log.Sugar.Fatal("config file path error")
		return
	}
	//debug = true
}

func Main() {
	var (
		err   error
		ready chan struct{} = make(chan struct{})
	)
	logic.InitNode(config.Cfg.Addr, config.Cfg.Group, config.Cfg.Rpc, config.Cfg.Peers)
	go node.RunRpcServer(config.Cfg.Addr+":"+config.Cfg.Rpc, logic.Self)

	if err = logic.InitEmbed(
		config.Cfg.Group,
		config.Cfg.Dir,
		config.Cfg.Addr,
		config.Cfg.Client,
		config.Cfg.Peer,
		config.Cfg.State,
		config.Cfg.Metrics,
		config.Cfg.Peers,
	); err != nil {
		log.Sugar.Error("init embed server failed. error: ", err)
	}

	go func() { // 运行etcd
		if err = logic.Device.Run(ready); err != nil {
			log.Sugar.Error("run embed server error. ", err)
			return
		}
	}()
	select {
	case <-ready: // 待etcd Ready 运行其他服务
		err = logic.Device.SetAuth(Username, Password)
		if err != nil {
			log.Sugar.Fatal("set auth embed server failed. error: ", err)
		}
		Server()
	}
}

func Server() {
	var (
		err error
	)
	// 客户端连接
	if err = logic.InitStore(config.Cfg.Addr, config.Cfg.Client, Username, Password, time.Second*3); err != nil {
		log.Sugar.Error("init store failed. error: ", err)
	}

	// 从etcd加载license
	if err := loadLic(); err != nil {
		log.Sugar.Error("init license failed. error: ", err)
	}
	logic.DefaultConf()
	logic.MemberConf(config.Cfg.Web)

	// 间隔1分钟更新授权
	go func() {
		ticker := time.Tick(1 * time.Minute) // 1分钟
		expr := cronexpr.MustParse("* * * * *")
		for range ticker {
			// 成员列表
			logic.MemberConf(config.Cfg.Web)
			now := time.Now()
			next := expr.Next(now)
			time.AfterFunc(next.Sub(now), func() {
				// 如果是主就更新授权
				if logic.Device.IsLeader() {
					log.Sugar.Infof("%s is Leader. ip:%s", logic.Self.Name, logic.Self.IP)
					if err := logic.ResetLicense(); err != nil {
						log.Sugar.Error("reset license failed. error: ", err)
					}
				}
			})
		}
	}()
	// 监听授权变化
	go logic.WatchLicense()
	go RunWebWithHttps(config.Cfg.Web)

	// 时间同步服务
	//if config.Cfg.Ntp {
	//	go func() {
	//		ntpd.Run()
	//	}()
	//}

	// 阻塞主进程
	<-make(chan struct{})
	//<- (chan int)(nil)
}

// 启动程序时加载授权
func loadLic() (err error) {
	var (
		cipher string
		lic    *model.License
	)
	if cipher, err = logic.GetLicense(); err != nil {
		log.Sugar.Error("get license failed. error: ", err)
	}

	if cipher == "" {
		lic = new(model.License)
	} else {
		lic, err = logic.Str2lic(cipher)
	}
	logic.StoreLic(lic)

	return
}

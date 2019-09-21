package odinmain

import (
	"../asset"
	"../config"
	"../log"
	"../logic"
	"../model"
	"../node"
	"fmt"
	"github.com/gorhill/cronexpr"
	//"../ntpd"
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
)

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
			log.Sugar.Error("restore assets failed. error: ", err.Error())
			_ = os.RemoveAll(filepath.Join(AssetPath, dir))
			continue
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(logo)
	//RestoreAsset()
}

func Main() {
	var (
		err   error
		ready chan struct{} = make(chan struct{})
	)
	go node.RunRpcServer(config.Cfg.Rpc, logic.Self)

	if err = logic.InitEmbed(); err != nil {
		log.Sugar.Error("init embed server failed. error: ", err.Error())
	}

	go func() { // 运行etcd
		if err = logic.Device.Run(ready); err != nil {
			log.Sugar.Error("run embed server error. ", err.Error())
			return
		}
	}()
	select {
	case <-ready: // 待etcd Ready 运行其他服务
		err = logic.Device.SetAuth(Username, Password)
		if err != nil {
			log.Sugar.Fatal("set auth embed server failed. error: ", err.Error())
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
		log.Sugar.Error("init store failed. error: ", err.Error())
	}
	// 生成机器码
	//if err = initHardware(); err != nil {
	//	log.Error("Failed to generate rank code when service starts. error: ", err.Error())
	//}

	// 从etcd加载license
	if err := loadLic(); err != nil {
		log.Sugar.Error("init license failed. error: ", err.Error())
	}

	// 间隔1分钟更新授权
	go func() {
		ticker := time.Tick(1 * time.Minute) // 1分钟
		expr := cronexpr.MustParse("* * * * *")
		for range ticker {
			now := time.Now()
			next := expr.Next(now)
			time.AfterFunc(next.Sub(now), func() {
				// 如果是主就更新授权
				if logic.Device.IsLeader() {
					log.Sugar.Infof("%s is Leader. ip:%s", logic.Self.Name, logic.Self.IP)
					if err := logic.ResetLicense(); err != nil {
						log.Sugar.Error("reset license failed. error: ", err.Error())
					}
				}
			})
		}
	}()
	// 监听授权变化
	go logic.WatchLicense()
	// 监听客户端
	//go func() {
	//	if err := EH.WatchClient(); err != nil {
	//		log.Error("watch CollectionCli error. error: ", err.Error())
	//	}
	//}()
	// web
	go RunWebWithHttps(config.Cfg.Web)

	//time.Sleep(3*time.Second)

	// 向etcd注册本节点
	//if err := initNodeStatus(); err != nil {
	//	log.Error("initialization self node failed. error: ", err.Error())
	//}

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
	err = logic.GetLicenseStatus()
	if err != nil {
		log.Sugar.Error("get license status failed. error: ", err.Error())
	}
	if cipher, err = logic.GetLicense(); err != nil {
		log.Sugar.Error("get license failed. error: ", err.Error())
	}

	lic, err = logic.Str2lic(cipher)
	logic.StoreLic(lic)

	return
}

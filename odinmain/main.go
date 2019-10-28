package odinmain

import (
	"fmt"
	"github.com/offer365/odin/asset"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/logic"
	"github.com/offer365/odin/proto"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"

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
	_assetPath string
	User             = "admin"
)

// 释放静态资源
func RestoreAsset() {
	// 解压 静态文件的位置
	if runtime.GOOS == "linux" {
		_assetPath = "/usr/share/.asset/.temp/"
	} else {
		_assetPath = "./"
	}
	// go get -u github.com/jteeuwen/go-bindata/...
	// 重新生成静态资源在项目的根目录下 go-bindata -o=asset/asset.go -pkg=asset html/... static/...
	dirs := []string{"html", "static"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets(_assetPath, dir); err != nil {
			log.Sugar.Error("restore assets failed. error: ", err)
			_ = os.RemoveAll(filepath.Join(_assetPath, dir))
			continue
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(logo)
	RestoreAsset()
}

func Main() {
	var (
		err   error
		ready = make(chan struct{})
	)

	if err = logic.InitEmbed(
		config.Cfg.Name,
		config.Cfg.Dir,
		config.Cfg.LocalClientAddr(),
		config.Cfg.LocalPeerAddr(),
		config.Cfg.State,
		config.Cfg.AllPeerAddr(),
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
	if err = logic.InitStore(config.Cfg.LocalClientAddr(), Username, Password, time.Second*3); err != nil {
		log.Sugar.Error("init store failed. error: ", err)
	}

	// 从etcd加载license
	if err := loadLic(); err != nil {
		log.Sugar.Error("init license failed. error: ", err)
	}

	// 间隔1分钟更新授权
	go func() {
		ticker := time.Tick(1 * time.Minute) // 1分钟
		//expr := cronexpr.MustParse("* * * * *")
		for range ticker {
			//now := time.Date()
			//next := expr.Next(now)
			//time.AfterFunc(next.Sub(now), func() {
			//time.AfterFunc(time.Second, func() {})
			// 如果是主就更新授权
			if logic.Device.IsLeader() {
				log.Sugar.Infof("%s is Leader. ip:%s", proto.Self.Attrs.Name, proto.Self.Attrs.Addr)
				if err := logic.ResetLicense(); err != nil {
					log.Sugar.Error("reset license failed. error: ", err)
				}
			}
		}
	}()
	// 监听授权变化
	go logic.WatchLicense()
	go Run(config.Cfg.LocalGRpcAddr())
	proto.AllNodeGRpcClient(config.Cfg.AllGRpcAddr())
	logic.DefaultConf()
	signalChan := make(chan os.Signal)
	done := make(chan struct{}, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, os.Kill)
	// 资源回收
	go func() {
		<-signalChan
		proto.ClientConns.Range(func(key, value interface{}) bool {
			cli, ok := value.(*grpc.ClientConn)
			if ok {
				cli.Close()
			}
			return true
		})
		gs.Stop()
		logic.Close()
		done <- struct{}{}
	}()
	// 阻塞主进程
	<-done
	//<-make(chan struct{})
	//<- (chan int)(nil)
}

// 启动程序时加载授权
func loadLic() (err error) {
	var (
		byt []byte
		lic    *logic.License
	)
	if byt, err = logic.GetLicense(); err != nil {
		log.Sugar.Error("get license failed. error: ", err)
	}

	if byt == nil || len(byt) == 0 {
		lic = new(logic.License)
	} else {
		lic, err = logic.Str2lic(string(byt))
	}
	logic.StoreLic(lic)
	return
}

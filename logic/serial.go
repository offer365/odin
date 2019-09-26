package logic

import (
	"context"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/node"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"sync"
	"time"
)

///////////////////////////////////////////////////////////////////////////////////////
// 序列号																			//
//////////////////////////////////////////////////////////////////////////////////////

// 从etcd获取序列号
func GetSerialNum() (info string, err error) {
	var (
		getResp *clientv3.GetResponse
	)
	info = "请重新获取序列号。"
	if getResp, err = store.Get(serialNumKey); err != nil {
		return
	}
	if len(getResp.Kvs) > 0 {
		info = string(getResp.Kvs[0].Value)
	}
	return
}

// 写入序列号
func PutSerialNum(val string) (err error) {
	if _, err = store.Put(serialNumKey, val); err != nil {
		log.Sugar.Error("put serial num failed. error: ", err.Error())
	}
	return nil
}

func ResetSerialNum() (code string, err error) {
	code, err = Serial.GenSerialNum(GetAllNodes())
	if err != nil {
		goto ERR
	}
	if err = PutSerialNum(code); err != nil {
		goto ERR
	}
	return
ERR:
	log.Sugar.Error("reset serial num failed. error: ", err.Error())
	return
}

func GetAllNodes() (nodes map[string]*node.Node) {
	var lock sync.Mutex
	var wait sync.WaitGroup
	//var  value atomic.Value
	nodes = make(map[string]*node.Node, 0)
	//value.Store(nodes)
	wait.Add(len(config.Cfg.Peers))
	for id, ip := range config.Cfg.Peers {
		go func(id int, ip string) {
			defer wait.Done()
			name := "odin" + strconv.Itoa(id)
			ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
			n, err := node.GetRemoteNode(ctx, name, ip, config.Cfg.Rpc)
			if err != nil {
				log.Sugar.Error("node rpc dial failed. error: ", err)
				return
			}
			lock.Lock()
			nodes[name] = n
			lock.Unlock()
			return
		}(id, ip)
	}
	wait.Wait()
	//sort.Slice(nodes, func(i, j int) bool {
	//	return nodes[i].Name < nodes[j].Name
	//})
	return
}

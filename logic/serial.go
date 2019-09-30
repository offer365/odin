package logic

import (
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/node"
	"go.etcd.io/etcd/clientv3"
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
	if getResp, err = store.Get(serialNumKey, true); err != nil {
		return
	}
	if len(getResp.Kvs) > 0 {
		info = string(getResp.Kvs[0].Value)
	}
	return
}

// 写入序列号
func PutSerialNum(val string) (err error) {
	if _, err = store.Put(serialNumKey, val, true); err != nil {
		log.Sugar.Error("put serial num failed. error: ", err)
	}
	return nil
}

func ResetSerialNum() (code string, err error) {
	nodes:=node.GetAllNodes(Self.Rpc,Self.Peers)
	code, err = Serial.GenSerialNum(nodes)
	if err != nil {
		goto ERR
	}
	if err = PutSerialNum(code); err != nil {
		goto ERR
	}
	return
ERR:
	log.Sugar.Error("reset serial num failed. error: ", err)
	return
}

func GetAllNode() map[string]*node.Node {
	return node.GetAllNodes(Self.Rpc,Self.Peers)
}


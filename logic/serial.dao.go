package logic

import (
	"context"
	"time"

	"github.com/offer365/odin/log"
	pb "github.com/offer365/odin/proto"
	"go.etcd.io/etcd/clientv3"
)

// 序列号

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

// 重置序列号
func ResetSerialNum() (code string, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	nodes := pb.GetAllNodes(ctx)
	code, err = Serial.Generate(nodes)
	if err != nil {
		log.Sugar.Error("generate serial num failed. error: ", err)
		return
	}
	if err = PutSerialNum(code); err != nil {
		log.Sugar.Error("put serial num failed. error: ", err)
		return
	}
	return
}

func GetAllNode() map[string]*pb.Node {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	return pb.GetAllNodes(ctx)
}

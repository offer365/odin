package logic

import (
	"encoding/json"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/model"
	"go.etcd.io/etcd/clientv3"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////////////
// 客户端实例	                                                                    //
/////////////////////////////////////////////////////////////////////////////////////

// 获取Client
func GetClient(key string) (cli *model.Cli, ok bool) {
	var (
		resp *clientv3.GetResponse
	)
	key = clientKeyPrefix + key
	resp, err := store.Get(key, false)
	if err != nil {
		log.Sugar.Error("get client failed. error: ", err)
		return
	}
	if len(resp.Kvs) > 0 {
		cli = new(model.Cli)
		err = json.Unmarshal(resp.Kvs[0].Value, cli)
		if err != nil {
			return nil, false
		}
		ok = true
	}
	return
}

// 获取所有Client
func GetAllClient(app string) (all map[string]string, err error) {
	var (
		getResp *clientv3.GetResponse
	)
	key := clientKeyPrefix + app
	if getResp, err = store.GetAll(key, false); err != nil {
		log.Sugar.Error("get all client failed. error: ", err)
		return
	}
	all = make(map[string]string, 0)
	for _, i := range getResp.Kvs {
		// TODO 是否字符串切分
		key := strings.Split(string(i.Key), clientKeyPrefix)[1]
		all[key] = string(i.Value)
	}
	return
}

// 获取所有Client个数
func ClientCount(app string) (count int64, err error) {
	var (
		resp *clientv3.GetResponse
	)
	key := clientKeyPrefix + app
	if resp, err = store.Count(key, true); err != nil {
		log.Sugar.Error("get all client failed. error: ", err)
		return
	}
	return resp.Count, err
}

// 写入Client
func PutClient(key string, cli *model.Cli) (lease int64, err error) {
	key = clientKeyPrefix + key
	// 10秒租期
	lg, err := store.Lease(key, 10)
	if err != nil {
		return
	}
	cli.Lease = int64(lg.ID)
	byt, err := json.Marshal(cli)
	if err != nil {
		return
	}
	if _, err = store.PutWithLease(key, string(byt), lg.ID, false); err != nil {
		return
	}
	lease = int64(lg.ID)
	return
}

// 删除Client
func DelClient(key string, leaseId int64) (err error) {
	key = clientKeyPrefix + key
	_, err = store.DelWithLease(key, leaseId, false)
	return
}

// 续租
func KeepAliveClient(key string, leaseId int64) (err error) {
	key = clientKeyPrefix + key
	_, err = store.KeepOnce(key, leaseId)
	return
}

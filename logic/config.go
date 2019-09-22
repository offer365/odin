package logic

import (
	"github.com/offer365/odin/log"
	"go.etcd.io/etcd/clientv3"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////////////
// 配置管理																			//
//////////////////////////////////////////////////////////////////////////////////////

// 获取配置
func GetConfig(key string) (val string, err error) {
	var (
		resp *clientv3.GetResponse
	)

	key = clientConfigKeyPrefix + key
	if resp, err = store.Get(key); err != nil {
		log.Sugar.Errorf("get %s config failed. error: %s", key, err.Error())
		return
	}
	if len(resp.Kvs) > 0 {
		val = string(resp.Kvs[0].Value)
	}
	return
}

// 获取所有配置
func GetAllConfig() (infos map[string]string, err error) {
	var (
		getResp *clientv3.GetResponse
	)
	if getResp, err = store.GetAll(clientConfigKeyPrefix); err != nil {
		log.Sugar.Error("get all config failed. error: ", err.Error())
		return
	}
	infos = make(map[string]string, 0)
	for _, i := range getResp.Kvs {
		// TODO 是否字符串切分
		key := strings.Split(string(i.Key), clientConfigKeyPrefix)[1]
		infos[key] = string(i.Value)
	}
	return
}

// 写入配置
func PutConfig(key, info string) (err error) {
	key = clientConfigKeyPrefix + key
	if _, err = store.Put(key, info); err != nil {
		log.Sugar.Error("put config failed. error: ", err.Error())
	}
	return
}

// 删除配置
func DelConfig(key string) (err error) {
	key = clientConfigKeyPrefix + key
	if _, err = store.Del(key); err != nil {
		log.Sugar.Error("del config failed. error: ", err.Error())
	}
	return
}

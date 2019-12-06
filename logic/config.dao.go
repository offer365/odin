package logic

import (
	"strings"

	"github.com/offer365/odin/log"
	"go.etcd.io/etcd/clientv3"
)

// 配置管理

// 获取配置
func GetConfig(key string) (val string, err error) {
	var (
		resp *clientv3.GetResponse
	)

	key = clientConfigKeyPrefix + key
	if resp, err = store.Get(key, true); err != nil {
		log.Sugar.Errorf("get %s config failed. error: %s", key, err)
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
	if getResp, err = store.GetAll(clientConfigKeyPrefix, true); err != nil {
		log.Sugar.Error("get all config failed. error: ", err)
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
	if _, err = store.Put(key, info, true); err != nil {
		log.Sugar.Error("put config failed. error: ", err)
	}
	return
}

// 删除配置
func DelConfig(key string) (err error) {
	_, ok := confWhiteList[key]
	if ok {
		return
	}
	key = clientConfigKeyPrefix + key
	if _, err = store.Del(key, true); err != nil {
		log.Sugar.Error("del config failed. error: ", err)
	}
	return
}

func DefaultConf() {
	var (
		df  string
		err error
	)

	if df, err = GetConfig(defaultKey); err != nil {
		log.Sugar.Error("get config default failed. error: ", err)
		return
	}
	if df == "" {
		if err = PutConfig(defaultKey, defaultKey); err != nil {
			log.Sugar.Error("put config default failed. error: ", err)
			return
		}
	}
}

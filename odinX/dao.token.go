package odinX

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/offer365/odin/log"
	"github.com/offer365/odin/utils"
	"go.etcd.io/etcd/clientv3"
)


// 下面的函数并不会直接存储App发送的token,而是对app,ID,val 分别进行md5或sha256 hash 后,再进行存储或取出。上层函数以同样的算法计算后，操作和比对。

func PutToken(app, id, val string) (err error) {
	key := Cfg.StoreTokenKey + Cfg.TokenHash([]byte(app)) + "/" + Cfg.TokenHash([]byte(id))
	val = Cfg.TokenHash([]byte(val))
	if _, err = store.Put(key, val, true); err != nil {
		log.Sugar.Error("put Token failed. error: ", err)
		return
	}
	return
}

func GetToken(app, id string) (token string, err error) {
	key := Cfg.StoreTokenKey + Cfg.TokenHash([]byte(app)) + "/" + Cfg.TokenHash([]byte(id))
	resp, err := store.Get(key, false)
	if err != nil {
		log.Sugar.Error("get client failed. error: ", err)
		return
	}
	if len(resp.Kvs) > 0 {
		token = string(resp.Kvs[0].Value)
	}
	return
}

// 删除Client
func DelToken(app, id string) (err error) {
	key := Cfg.StoreTokenKey + Cfg.TokenHash([]byte(app)) + "/" + Cfg.TokenHash([]byte(id))
	_, err = store.Del(key, true)
	return
}

func CountTokenWithApp(app string) (count int64, err error) {
	var (
		resp *clientv3.GetResponse
	)
	key := Cfg.StoreTokenKey + Cfg.TokenHash([]byte(app)) + "/"
	if resp, err = store.Count(key, true); err != nil {
		log.Sugar.Error("get count Token failed. error: ", err)
		return
	}
	return resp.Count, err
}

func CountTokenWithID(app, id string) (count int64, err error) {
	var (
		resp *clientv3.GetResponse
	)
	key := Cfg.StoreTokenKey + Cfg.TokenHash([]byte(app)) + "/" + Cfg.TokenHash([]byte(id))
	if resp, err = store.Count(key, true); err != nil {
		log.Sugar.Error("get count Token failed. error: ", err)
		return
	}
	return resp.Count, err
}

// 检查Token 是否正确或者不存在，是否可以注册
func GetTokenAndChk(app, id, token string) (exist, register bool) {
	var (
		err    error
		result string
		num    int64
	)
	if app == "" || id == "" || token == "" {
		return
	}
	// 尝试获取 Token
	if result, err = GetToken(app, id); err != nil {
		return
	}

	// 如果 一致 直接返回 存在
	if result == Cfg.TokenHash([]byte(token)) {
		exist = true
		return
	}
	// 如果该 app 的token 不正确 说明客户端token 改变，直接返回，需要解绑
	if result != "" {
		return
	}
	// 如果没有获取到即 token == "" ，则判断是否可以注册token
	if num, err = CountTokenWithApp(app); err != nil {
		return
	}
	// 如果大于授权则不添加
	register = LoadLic().ChkInstance(app, num)
	return
}

func Untied(app, id, code string) (err error) {
	var (
		byt []byte
	)
	if byt, err = base64.StdEncoding.DecodeString(code); err != nil {
		return
	}

	// 解密 解密后的内容是 一个map k,v经过Cfg.TokenHash计算
	if byt, err = Cfg.UntiedDecrypt(byt); err != nil {
		return
	}
	untie := new(UntiedCode)
	if err = json.Unmarshal(byt, untie); err != nil {
		return
	}
	// 计算 比对 解密后的值
	sha256Key := Cfg.TokenHash([]byte(app))
	sha256Val := Cfg.TokenHash([]byte(id))

	if untie.Key != sha256Key || untie.Value != sha256Val || utils.Abs(time.Now().Unix()-untie.Date) > 3600*6 {
		err = errors.New("app or id or date error ")
		return
	}
	return DelToken(app, id)
}

type UntiedCode struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Date  int64  `json:"date"`
}

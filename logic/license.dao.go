package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/offer365/endecrypt"
	"github.com/offer365/odin/log"
	pb "github.com/offer365/odin/proto"
	"go.etcd.io/etcd/clientv3"
	"sync/atomic"
	"time"
)

var atomicLic atomic.Value

// license

// 获取license
func GetLicense() (byt []byte, err error) {
	var (
		resp *clientv3.GetResponse
	)

	if resp, err = store.Get(licenseKey, true); err != nil {
		log.Sugar.Error("get license failed. error: ", err)
		return
	}
	if len(resp.Kvs) > 0 {
		byt = resp.Kvs[0].Value
	}
	return
}

// 刷新license
func PutLicense(val string) (err error) {
	if _, err = store.Put(licenseKey, val, true); err != nil {
		log.Sugar.Error("put license failed. error: ", err)
		return
	}
	return
}

// 删除license
func DelLicense() (err error) {
	if _, err = store.Del(licenseKey, true); err != nil {
		log.Sugar.Error("del license failed. error: ", err)
		return
	}
	return
}

// 反序列化license
func Str2lic(cipher string) (lic *License, err error) {
	var (
		byt []byte
	)
	if byt, err = base64.StdEncoding.DecodeString(cipher); err != nil {
		log.Sugar.Error("base64 license byte failed. error: ", err)
		return
	}
	if byt == nil || len(byt) == 0 {
		return
	}
	lic = new(License)
	if byt, err = endecrypt.Decrypt(endecrypt.Pub1AesRsa2048, byt); err != nil {
		log.Sugar.Errorf("decrypt license failed. error: ", err)
		return
	}
	if err = json.Unmarshal(byt, lic); err != nil {
		log.Sugar.Error("unmarshal failed. error: ", err)
	}
	return
}

// 序列化license
func lic2Str(lic interface{}) (cipher string, err error) {
	var (
		byt []byte
	)
	if byt, err = json.Marshal(lic); err != nil {
		log.Sugar.Errorf("pack license failed. error: ", err)
		return
	}

	if byt, err = endecrypt.Encrypt(endecrypt.Pri1AesRsa2048, byt); err != nil {
		log.Sugar.Error("encrypt failed. error: ", err)
		return
	}
	return base64.StdEncoding.EncodeToString(byt),err
}

// 重置license
func ResetLicense() (err error) {
	var (
		byt []byte
		cipher string
		lic *License
	)
	if byt, err = GetLicense(); err != nil {
		log.Sugar.Error("get lic failed. error: ", err)
		return
	}
	if byt == nil || len(byt) == 0 {
		err = errors.New("not found license")
		log.Sugar.Error(err)
		return
	}
	if lic, err = Str2lic(string(byt)); err != nil {
		log.Sugar.Error("in reset lic,get lic failed. error: ", err)
		return
	}
	if lic == nil || lic.Lid == "" {
		err = errors.New("license instance is nil")
		log.Sugar.Error(err)
		return
	}

	now := time.Now().Unix()
	num := (now - lic.Generate) / 60
	// 这里限制了 LifeCycle 只能不断的增大
	if num > lic.LifeCycle {
		//atomic.StoreInt64(&(lic.LifeCycle),num)
		lic.LifeCycle = num
	} else {
		//atomic.AddInt64(&(lic.LifeCycle),1)
		lic.LifeCycle += 1
	}

	// 这里限制了 UpdateTime 只能不断的增大
	if now > lic.Update {
		//atomic.StoreInt64(&(lic.Update),now)
		lic.Update = now
	} else {
		//atomic.AddInt64(&(lic.Update), 60)
		lic.Update += 60
	}
	if cipher, err = lic2Str(lic); err != nil {
		log.Sugar.Error("reset lic failed. error: ", err)
		return
	}
	if err = PutLicense(cipher); err != nil {
		log.Sugar.Error("reset lic failed. error: ", err)
	}
	StoreLic(lic)
	return
}

// 监听license
func WatchLicense() {
	putFunc := func(event *clientv3.Event) error {
		if !Device.IsLeader() {
			lic, err := Str2lic(string(event.Kv.Value))
			if err == nil {
				StoreLic(lic)
			}
			return err
		}
		return nil
	}
	delFunc := func(event *clientv3.Event) error {
		lic := LoadLic()
		lic.Apps = make(map[string]*App, 0)
		StoreLic(lic)
		return nil
	}

	ctx, _ := context.WithCancel(context.Background())
	store.Watch(ctx, licenseKey, putFunc, delFunc)
}

// 检查授权码是否合法
func ChkLicense(cipher string) (lic *License, ok bool, msg string) {
	var (
		err error
	)

	if lic, err = Str2lic(cipher); err != nil {
		msg = "未能正确解析授权码。"
		return
	}
	// 当前机器是否在授权中
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	nodes := pb.GetAllNodes(ctx)
	if len(nodes) != len(lic.Devices) {
		msg = "节点数量不一致。"
		return
	}
	hw, exist := lic.Devices[pb.Self.Attrs.Name]
	if !exist {
		msg = "未在授权中找到本机id。"
		return
	}
	if hw != pb.Self.Attrs.Hwmd5 {
		msg = "绑定硬件信息错误。"
		return
	}
	// 检查 时间是否正确 如果授权时间大于当前服务器时间，说明当前服务器时间慢
	now := time.Now().Unix()
	if lic.Update >= now || lic.Generate >= now {
		msg = "当前服务器时间不正确。"
		return
	}
	// 授权码6个小时有效期

	if (now-lic.Update) > 60*60*6 || (now-lic.Generate) > 60*60*6 {
		msg = "授权码超时失效。"
		return
	}
	// 如果 新到授权码的id 与当前授权码的id一致说明是重复授权
	if lic.Lid == LoadLic().Lid {
		msg = "重复的授权码。"
		return
	}
	// 唯一有用的值是sid
	if Serial.Sid != lic.Sid {
		msg = "序列号与授权码不一致。"
		return
	}
	ok = true
	msg = "授权码正确"
	return
}

func LoadLic() *License {
	val := atomicLic.Load()
	lic, ok := val.(*License)
	if ok {
		return lic
	}
	return new(License)
}

func StoreLic(lic *License) {
	atomicLic.Store(lic)
}

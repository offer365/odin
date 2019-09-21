package logic

import (
	"../log"
	"../model"
	"encoding/json"
	"github.com/offer365/endecrypt"
	"github.com/offer365/endecrypt/endeaesrsa"
	"go.etcd.io/etcd/clientv3"
	"sync/atomic"
	"time"
)

var atomicLic atomic.Value

///////////////////////////////////////////////////////////////////////////////////////
// license																			//
//////////////////////////////////////////////////////////////////////////////////////

// 从etcd获取license
func GetLicense() (cipher string, err error) {
	var (
		resp *clientv3.GetResponse
	)

	if resp, err = store.Get(licenseKey); err != nil {
		log.Sugar.Error("get license failed. error: ", err)
		return
	}
	if len(resp.Kvs) > 0 {
		cipher = string(resp.Kvs[0].Value)
	}
	return
}

// 刷新license
func PutLicense(val string) (err error) {
	if _, err = store.Put(licenseKey, val); err != nil {
		log.Sugar.Error("put license failed. error: ", err)
		return
	}
	IsAuth = true
	// 将授权状态变为 true
	if _, err = store.Put(licenseStatusKEY, "true"); err != nil {
		log.Sugar.Error("get license status failed. error: ", err)
		return
	}
	return
}

// 删除license
func DelLicense() (err error) {
	if _, err = store.Del(licenseKey); err != nil {
		log.Sugar.Error("del license failed. error: ", err)
		return
	}
	IsAuth = false
	// 将授权状态变为 false
	if _, err = store.Put(licenseStatusKEY, "false"); err != nil {
		log.Sugar.Error("get license status failed. error: ", err)
		return
	}
	return
}

// 反序列化license
func Str2lic(text string) (lic *model.License, err error) {
	lic = new(model.License)
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Error("str to license failed. error: ", err)
			return
		}
	}()
	if text, err = endeaesrsa.PubDecrypt(text, endecrypt.PubkeyServer2048, endecrypt.AesKeyServer2); err != nil {
		log.Sugar.Errorf("decrypt license failed. error: ", err.Error())
		return
	}
	//lic.Devices = make(map[string]string)
	//lic.APPs=make(map[string]*model.APP)
	if err = json.Unmarshal([]byte(text), lic); err != nil {
		log.Sugar.Error("unmarshal failed. error: ", err)
	}
	return
}

// 序列化license
func lic2str(lic interface{}) (text string, err error) {
	var (
		src []byte
	)
	if src, err = json.Marshal(lic); err != nil {
		log.Sugar.Errorf("pack license failed. error: ", err.Error())
		return
	}
	if text, err = endeaesrsa.PriEncrypt(src, endecrypt.PirkeyServer2048, endecrypt.AesKeyServer2); err != nil {
		log.Sugar.Error("encrypt failed. error: ", err)
	}
	return
}

// 重置license
func ResetLicense() (err error) {
	var (
		text string
	)
	if !IsAuth {
		return nil
	}
	now := time.Now()
	lic := LoadLic()
	num := (now.Unix() - lic.GenerationTime.Unix()) / 60
	if num > lic.LifeCycle {
		lic.LifeCycle = num
	} else {
		lic.LifeCycle += 1
	}

	if now.After(lic.UpdateTime) {
		lic.UpdateTime = now
	} else {
		lic.UpdateTime.Add(60 * time.Second)
	}
	if text, err = lic2str(lic); err != nil {
		log.Sugar.Error("reset lic failed. error: ", err.Error())
		return
	}
	if err = PutLicense(text); err != nil {
		log.Sugar.Error("reset lic failed. error: ", err.Error())
	}
	StoreLic(lic)
	return
}

// 监听license
func WatchLicense() {
	putFunc := func(event *clientv3.Event) error {
		IsAuth = true
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
		IsAuth = false
		lic := LoadLic()
		lic.APPs = make(map[string]*model.APP, 0)
		StoreLic(lic)
		return nil
	}

	store.Watch(licenseKey, putFunc, delFunc)
}

// 检查授权码是否合法
func ChkLicense(cipher string) (lic *model.License, ok bool, msg string) {
	var (
		err error
	)

	if lic, err = Str2lic(cipher); err != nil {
		msg = "未能正确解析授权码。"
		return
	}
	// 当前机器是否在授权中
	allHw := GetAllNodes()
	if len(allHw) != len(lic.Devices) {
		msg = "节点数量不一致。"
		return
	}
	hw, exist := lic.Devices[Self.Name]
	if !exist {
		msg = "未在授权中找到本机id。"
		return
	}
	if hw != Self.HwMd5 {
		msg = "绑定硬件信息错误。"
		return
	}
	// 检查 时间是否正确 如果授权时间大于当前服务器时间，说明当前服务器时间慢
	if lic.UpdateTime.After(time.Now()) {
		msg = "当前服务器时间不正确。"
		return
	}
	// 授权码6个小时有效期

	if (time.Now().Unix() - lic.UpdateTime.Unix()) > 60*60*6 {
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

// 获取授权状态
func GetLicenseStatus() (err error) {
	var (
		resp *clientv3.GetResponse
	)
	if resp, err = store.Get(licenseStatusKEY); err != nil {
		IsAuth = false
		return
	}
	if len(resp.Kvs) > 0 {
		if string(resp.Kvs[0].Value) == "true" {
			IsAuth = true
		}
	}

	if resp, err = store.Count(licenseKey); err != nil {
		IsAuth = false
		return
	}
	if resp.Count == 0 {
		IsAuth = false
	}
	return
}

func LoadLic() *model.License {
	val := atomicLic.Load()
	lic, ok := val.(*model.License)
	if ok {
		return lic
	}
	return new(model.License)
}

func StoreLic(lic *model.License) {
	atomicLic.Store(lic)
}

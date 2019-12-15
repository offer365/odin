package odinX

import (
	"time"

	"go.etcd.io/etcd/clientv3"
)

// 生成注销码
func GenClearLicense() (text string, err error) {
	byt, err := GetLicense()
	if err != nil {
		return
	}
	if err = DelLicense(); err != nil {
		return
	}
	lic := LoadLic()
	lic.Apps = make(map[string]*App, 0)
	StoreLic(lic)
	obj := &Clear{
		Lic:    lic,
		Cipher: string(byt),
		Date:   time.Now().Unix(),
	}
	text, err = lic2Str(obj)
	if err != nil {
		return
	}
	return text, PutClearLicense(text)
}

// 获取注销码
func GetClearLicense() (info string, err error) {
	var (
		getResp *clientv3.GetResponse
	)
	if getResp, err = store.Get(Cfg.StoreClearLicenseKey, true); err != nil {
		return
	}
	if len(getResp.Kvs) > 0 {
		info = string(getResp.Kvs[0].Value)
	}
	return
}

// 写入注销码
func PutClearLicense(info string) (err error) {
	_, err = store.Put(Cfg.StoreClearLicenseKey, info, true)
	return
}

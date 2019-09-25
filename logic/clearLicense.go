package logic

import (
	"github.com/offer365/odin/model"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 生成注销码
func GenClearLicense() (text string, err error) {
	if err = DelLicense(); err != nil {
		return
	}
	lic := LoadLic()
	lic.APPs = make(map[string]*model.APP, 0)
	StoreLic(lic)
	text, err = GetLicense()
	if err != nil {
		return
	}
	obj := make(map[string]interface{})
	obj["real_time_license"] = lic
	obj["license"] = text
	obj["date"] = time.Now().Format("2006-01-02 15:04:05")
	text, err = lic2str(obj)
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
	if getResp, err = store.Get(clearLicenseKey); err != nil {
		return
	}
	if len(getResp.Kvs) > 0 {
		info = string(getResp.Kvs[0].Value)
	}
	return
}

// 写入注销码
func PutClearLicense(info string) (err error) {
	_, err = store.Put(clearLicenseKey, info)
	return
}

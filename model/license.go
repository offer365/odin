package model

import (
	"reflect"
	"strconv"
	"sync"
	"time"
)

var fml *fLicense

func init() {
	fml = new(fLicense)
}

type fLicense struct {
	UpdateTime string  `json:"update_time" title:"更新时间"` //当前时间 最后一次授权更新时间
	LifeCycle  int64   `json:"life_cycle" title:"生存周期"`  // 当前生存周期
	Apps       []*fApp `json:"apps"`
}

type fApp struct {
	Title string `json:"title"`
	Data  []item `json:"data"`
}

type item struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

// 授权码
type License struct {
	Lid            string            `json:"lid"`                      // 授权码唯一uuid,用来甄别是否重复授权。
	Sid            string            `json:"sid"`                      // 机器码的id, lid与sid 一一对应
	Devices        map[string]string `json:"devices"`                  // 节点id与 硬件信息md5
	GenerationTime time.Time         `json:"generation_time"`          // 授权生成时间
	UpdateTime     time.Time         `json:"update_time" title:"更新时间"` //当前时间 最后一次授权更新时间
	LifeCycle      int64             `json:"life_cycle" title:"生存周期"`  // 当前生存周期
	APPs           map[string]*APP   `json:"apps"  title:"产品"`         //key:app英文名请求中url标识字段
}

type APP struct {
	Url          string           `json:"url"`
	Name         string           `json:"name" title:"服务"`
	Attr         map[string]int64 `json:"attr"`                          // 自定义内容
	Instance     int              `json:"instance" title:"最大实例"`         // 实例
	ExpireTime   time.Time        `json:"expire_time" title:"到期时间"`      // 授权到期的时间戳
	MaxLifeCycle int64            `json:"max_life_cycle" title:"最大生存周期"` // 最大生存周期 (授权到期时间-生成授权时间)/周期时间60s
	rv           reflect.Value
	rt           reflect.Type
	mu           sync.RWMutex
}

// 传入一个app名，检查改产品是否到期
func (lic *License) CheckTime(application string) bool {
	app, ok := lic.APPs[application]
	if ok {
		now := time.Now()
		return now.Before(app.ExpireTime) && lic.LifeCycle < app.MaxLifeCycle && (now.Unix()-lic.GenerationTime.Unix())/60 < app.MaxLifeCycle && lic.UpdateTime.Before(app.ExpireTime)
	}

	return false
}

func (lic *License) NodeNums() int {
	return len(lic.Devices)
}

func (lic *License) Format() *fLicense {
	fml.UpdateTime = lic.UpdateTime.Format("2006-01-02 15:04:05")
	fml.LifeCycle = lic.LifeCycle
	fml.Apps = make([]*fApp, 0) //清空数组

	for _, app := range lic.APPs {
		app.reflect()
		fapp := new(fApp)
		fapp.Title = app.Name
		fapp.Data = make([]item, 0)
		//fapp.Data = append(fapp.Data, app.fieldName())
		for _, i := range app.fieldAttr() {
			fapp.Data = append(fapp.Data, i)
		}
		fapp.Data = append(fapp.Data, app.fieldInstance())
		fapp.Data = append(fapp.Data, app.fieldExpireTime())
		fapp.Data = append(fapp.Data, app.fieldMaxLifeCycle())
		fml.Apps = append(fml.Apps, fapp)
	}
	return fml
}

func (a *APP) reflect() {
	a.rv = reflect.ValueOf(*a)
	a.rt = a.rv.Type()
}

func (a *APP) fieldName() item {
	var (
		title string
		value string
	)
	name, ok := a.rt.FieldByName("name")
	if ok {
		title = name.Tag.Get("title")
		value = a.rv.FieldByName("name").Interface().(string)
	}
	return item{
		Title: title,
		Value: value,
	}
}

func (a *APP) fieldAttr() (items []item) {
	var (
		value map[string]int64
	)
	value, ok := a.rv.FieldByName("Attr").Interface().(map[string]int64)
	if ok {
		for k, v := range value {
			items = append(items, item{Title: k, Value: strconv.Itoa(int(v))})
		}
	}
	return
}

func (a *APP) fieldInstance() item {
	var (
		title string
		value int
	)
	instance, ok := a.rt.FieldByName("Instance")
	if ok {
		title = instance.Tag.Get("title")
		value = a.rv.FieldByName("Instance").Interface().(int)
	}
	return item{
		Title: title,
		Value: strconv.Itoa(value),
	}
}

func (a *APP) fieldExpireTime() item {
	var (
		title string
		value time.Time
	)
	instance, ok := a.rt.FieldByName("ExpireTime")
	if ok {
		title = instance.Tag.Get("title")
		value = a.rv.FieldByName("ExpireTime").Interface().(time.Time)
	}
	return item{
		Title: title,
		Value: value.Format("2006-01-02 15:04:05"),
	}
}

func (a *APP) fieldMaxLifeCycle() item {
	var (
		title string
		value int64
	)
	instance, ok := a.rt.FieldByName("MaxLifeCycle")
	if ok {
		title = instance.Tag.Get("title")
		value = a.rv.FieldByName("MaxLifeCycle").Interface().(int64)
	}
	return item{
		Title: title,
		Value: strconv.Itoa(int(value)),
	}
}

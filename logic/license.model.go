package logic

import (
	"reflect"
	"strconv"
	"time"
)

var fml *fLicense

func init() {
	fml = new(fLicense)
}

type fLicense struct {
	UpdateTime string  `json:"update_time" title:"更新时间"` // 当前时间 最后一次授权更新时间
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
	Lid       string            `json:"lid"`                    // 授权码唯一uuid,用来甄别是否重复授权。
	Sid       string            `json:"sid"`                    // 机器码的id, lid与sid 一一对应
	Devices   map[string]string `json:"devices"`                // 节点id与 硬件信息md5
	Generate  int64             `json:"generate"`               // 授权生成时间
	Update    int64             `json:"update" title:"更新时间"`    // 当前时间 最后一次授权更新时间
	LifeCycle int64             `json:"lifeCycle" title:"生存周期"` // 当前生存周期
	Apps      map[string]*App   `json:"apps"  title:"产品"`       // map[key]*App key=App.key
}

type App struct {
	Key          string  `json:"key"`
	Name         string  `json:"name" title:"服务"`
	Attrs        []*Attr `json:"attrs"`                       // 自定义内容
	Instance     int64   `json:"instance" title:"最大实例"`       // 实例
	Expire       int64   `json:"expire" title:"到期时间"`         // 授权到期的时间戳
	MaxLifeCycle int64   `json:"maxLifeCycle" title:"最大生存周期"` // 最大生存周期 (授权到期时间-生成授权时间)/周期时间60s
	rv           reflect.Value
	rt           reflect.Type
}

type Attr struct {
	Name  string
	Key   string
	Value int64
}

// 传入一个app名，检查改产品是否到期
func (lic *License) CheckTime(application string) bool {
	app, ok := lic.Apps[application]
	if ok {
		now := time.Now().Unix()
		return now < app.Expire && lic.Update < app.Expire && lic.LifeCycle < app.MaxLifeCycle && (now-lic.Generate)/60 < app.MaxLifeCycle
	}
	return false
}

func (lic *License) ChkInstance(application string, num int64) bool {
	app, ok := lic.Apps[application]
	if ok {
		return app.Instance > num
	}
	return false
}

func (lic *License) Format() *fLicense {
	fml.UpdateTime = time.Unix(lic.Update, 0).Format("2006-01-02 15:04:05")
	fml.LifeCycle = lic.LifeCycle
	fml.Apps = make([]*fApp, 0) // 清空数组

	for _, app := range lic.Apps {
		app.reflect()
		fapp := new(fApp)
		fapp.Title = app.Name
		fapp.Data = make([]item, 0)
		// fapp.Data = append(fapp.Data, App.fieldName())
		for _, i := range app.fieldAttrs() {
			fapp.Data = append(fapp.Data, i)
		}
		fapp.Data = append(fapp.Data, app.fieldInstance())
		fapp.Data = append(fapp.Data, app.fieldExpireTime())
		fapp.Data = append(fapp.Data, app.fieldMaxLifeCycle())
		fml.Apps = append(fml.Apps, fapp)
	}
	return fml
}

func (a *App) reflect() {
	a.rv = reflect.ValueOf(*a)
	a.rt = a.rv.Type()
}

func (a *App) fieldName() item {
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

func (a *App) fieldAttrs() (items []item) {
	var (
		attrs []*Attr
	)
	attrs, ok := a.rv.FieldByName("Attrs").Interface().([]*Attr)
	if ok {
		for _, attr := range attrs {
			items = append(items, item{Title: attr.Name, Value: strconv.Itoa(int(attr.Value))})
		}
	}
	return
}

func (a *App) fieldInstance() item {
	var (
		title string
		value int64
	)
	instance, ok := a.rt.FieldByName("Instance")
	if ok {
		title = instance.Tag.Get("title")
		value = a.rv.FieldByName("Instance").Interface().(int64)
	}
	return item{
		Title: title,
		Value: strconv.Itoa(int(value)),
	}
}

func (a *App) fieldExpireTime() item {
	var (
		title string
		value int64
	)
	instance, ok := a.rt.FieldByName("Expire")
	if ok {
		title = instance.Tag.Get("title")
		value = a.rv.FieldByName("Expire").Interface().(int64)
	}
	return item{
		Title: title,
		Value: time.Unix(value, 0).Format("2006-01-02 15:04:05"),
	}
}

func (a *App) fieldMaxLifeCycle() item {
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

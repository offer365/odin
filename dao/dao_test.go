package dao

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"testing"
	"time"
)

var store Store

func TestMd5(t *testing.T) {
	h := md5.New()
	h.Write([]byte("aa"))
	a := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(a)
}


func init() {
	store=NewStore()
	err:=store.Init(context.Background(),
		WithHost("10.0.0.200"),
		WithPort("12379"),
		WithUsername("root"),
		WithPassword("613f#8d164df4ACPF49@93a510df49!66f98b*d6"),
	WithTimeout(time.Millisecond*2000),
	)
	fmt.Println(err)
}

func TestEtcdStore_Count(t *testing.T) {
	resp,err:=store.Count("/odin",false)
	fmt.Println(err)
	fmt.Println(resp.Count)
}

func TestEtcdStore_Get(t *testing.T) {
	resp,err:=store.Get("/odin/license",true)
	fmt.Println(err)
	for _,kv:=range resp.Kvs{
		fmt.Println(string(kv.Value))
	}
}

func TestEtcdStore_GetAll(t *testing.T) {
	resp,err:=store.GetAll("/odin",true)
	fmt.Println(err)
	for _,kv:=range resp.Kvs{
		fmt.Println(kv.String())
	}
}

func TestEtcdStore_Put(t *testing.T) {
	_,err:=store.Put("/odin/test","hehe",true)
	fmt.Println(err)
	resp,err:=store.Get("/odin/test",true)
	fmt.Println(err)
	for _,kv:=range resp.Kvs{
		fmt.Println(string(kv.Value))
	}

}

func TestEtcdStore_Del(t *testing.T) {
	_,err:=store.Del("/odin/test",true)
	fmt.Println(err)
	resp,err:=store.Get("/odin/test",true)
	fmt.Println(err)
	for _,kv:=range resp.Kvs{
		fmt.Println(string(kv.Value))
	}
}

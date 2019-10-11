package dao

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"testing"
	"time"
)

var store Store

func init() {
	store = NewStore()
	err := store.Init(context.Background(),
		WithHost("10.0.0.200"),
		WithPort("12379"),
		WithUsername("root"),
		WithPassword("613f#8d164df4ACPF49@93a510df49!66f98b*d6"),
		WithTimeout(time.Millisecond*2000),
	)
	fmt.Println(err)
}

func TestEtcdStore_Count(t *testing.T) {
	resp, err := store.Count("/odin", false)
	fmt.Println(err)
	fmt.Println(resp.Count)
}

func TestEtcdStore_Get(t *testing.T) {
	resp, err := store.Get("/odin/license", true)
	fmt.Println(err)
	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Value))
	}
}

func TestEtcdStore_GetAll(t *testing.T) {
	resp, err := store.GetAll("/odin", true)
	fmt.Println(err)
	for _, kv := range resp.Kvs {
		fmt.Println(kv.String())
	}
}

func TestEtcdStore_Put(t *testing.T) {
	_, err := store.Put("/odin/test", "hehe", true)
	fmt.Println(err)
	resp, err := store.Get("/odin/test", true)
	fmt.Println(err)
	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Value))
	}

}

func TestEtcdStore_Del(t *testing.T) {
	_, err := store.Del("/odin/test", true)
	fmt.Println(err)
	resp, err := store.Get("/odin/test", true)
	fmt.Println(err)
	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Value))
	}
}

// 第1步获取租约id: store.Lease(key,5)
// 第2步写入值: store.PutWithLease(key,val,resp.ID,false)
func TestEtcdStore_PutWithLease(t *testing.T) {
	key := "/aa/bb"
	// 10秒租期
	sl, err := store.Lease(key, 10)
	if err != nil {
		return
	}
	if _, err = store.PutWithLease(key, "cccc", sl.ID, false); err != nil {
		return
	}
	fmt.Println(sl.ID)
	for i := 0; i <= 11; i++ {
		time.Sleep(time.Second)
		result, err := store.Get(key, false)
		fmt.Println(err)
		for _, kv := range result.Kvs {
			fmt.Println(kv.String())
			fmt.Println(string(kv.Value))
		}
	}
}

// 第1步获取租约id: store.Lease(key,5)
// 第2步写入值: store.PutWithLease(key,val,resp.ID,false)
func TestEtcdStore_DelWithLease(t *testing.T) {
	key := "/aa/bb"
	// 10秒租期
	sl, err := store.Lease(key, 10)
	if err != nil {
		return
	}
	if _, err = store.PutWithLease(key, "test", sl.ID, false); err != nil {
		return
	}
	fmt.Println(sl.ID)
	time.Sleep(time.Second)
	result, err := store.Get(key, false)
	fmt.Println(err)
	for _, kv := range result.Kvs {
		fmt.Println(kv.String())
		fmt.Println(string(kv.Value))
	}
	store.DelWithLease(key, sl.ID, false)
	result, err = store.Get(key, false)
	fmt.Println(err, result.Count)
	for _, kv := range result.Kvs {
		fmt.Println(kv.String())
		fmt.Println(string(kv.Value))
	}
}

// 第1步获取租约id: store.Lease(key,5)
// 第2步写入值: store.PutWithLease(key,val,resp.ID,false)
// 第3步keepalive: store.KeepAlive(ctx,resp.ID)
func TestEtcdStore_KeepAlive(t *testing.T) {
	key := "/aaa/bbb"
	resp, _ := store.Lease(key, 5)
	store.PutWithLease(key, "test", resp.ID, false)
	ctx, cancel := context.WithCancel(context.Background())
	store.KeepAlive(ctx, resp.ID)
	go func() {
		for i := 0; i <= 50; i++ {
			time.Sleep(time.Second)
			result, _ := store.Get(key, false)
			for _, kv := range result.Kvs {
				fmt.Println(kv.String())
				fmt.Println(string(kv.Value))
			}
			fmt.Println(time.Now().String())
		}
	}()
	time.Sleep(time.Second * 15)
	cancel()
	time.Sleep(time.Second * 15)
}

// 第1步获取租约id: store.Lease(key,5)
// 第2步写入值: store.PutWithLease(key,val,resp.ID,false)
// 第3步keepOnce: store.KeepAlive(ctx,resp.ID)
func TestEtcdStore_KeepOnce(t *testing.T) {
	key := "/aaaaaa/bbbb"
	resp, _ := store.Lease(key, 5)
	store.PutWithLease(key, "test", resp.ID, false)
	go func() {
		for i := 0; i <= 10; i++ {
			store.KeepOnce(resp.ID)
			time.Sleep(3 * time.Second)
		}
	}()

	for i := 0; i <= 50; i++ {
		time.Sleep(time.Second)
		result, _ := store.Get(key, false)
		for _, kv := range result.Kvs {
			fmt.Println(kv.String())
			fmt.Println(string(kv.Value))
		}
		fmt.Println(time.Now().String())
	}
}

func TestEtcdStore_Watch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	putF := func(event *clientv3.Event) error {
		fmt.Println(string(event.Kv.Value))
		return nil
	}
	delF := func(event *clientv3.Event) error {
		fmt.Println(string(event.Kv.Key))
		return nil
	}
	go store.Watch(ctx, "/abcd", putF, delF)
	for i := 0; i < 10; i++ {
		store.Put("/abcd", strconv.Itoa(i+100), true)
		time.Sleep(time.Second)
	}
	store.Del("/abcd", true)
	cancel()
	for i := 0; i < 10; i++ {
		store.Put("/abcd", strconv.Itoa(i+100), true)
		time.Sleep(time.Second)
	}
}

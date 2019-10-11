package dao

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/offer365/odin/log"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

// etcd 客户端
type etcdStore struct {
	options *Options
	config  clientv3.Config
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
	timeout time.Duration
}

func (es *etcdStore) Init(ctx context.Context, opts ...Option) (err error) {
	es.options = new(Options)
	for _, opt := range opts {
		opt(es.options)
	}

	es.config = clientv3.Config{
		Endpoints:   []string{"http://" + es.options.Host + ":" + es.options.Port},
		DialTimeout: es.options.Timeout,
		Username:    es.options.Username,
		Password:    es.options.Password,
	}
	es.timeout = es.options.Timeout
	if es.client, err = clientv3.New(es.config); err != nil {
		log.Sugar.Error("create client failed. error: ", err)
		return
	}
	es.kv = clientv3.NewKV(es.client)
	es.lease = clientv3.NewLease(es.client)
	es.watcher = clientv3.NewWatcher(es.client)
	return
}

func (es *etcdStore) Get(key string, lock bool) (resp *clientv3.GetResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Get(ctx, key)
}

func (es *etcdStore) GetAll(prefix string, lock bool) (resp *clientv3.GetResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(prefix))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Get(ctx, prefix, clientv3.WithPrefix())
}

func (es *etcdStore) Count(prefix string, lock bool) (resp *clientv3.GetResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(prefix))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
}

func (es *etcdStore) Put(key, val string, lock bool) (resp *clientv3.PutResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Put(ctx, key, val)
}

func (es *etcdStore) Del(key string, lock bool) (resp *clientv3.DeleteResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Delete(ctx, key)
}

func (es *etcdStore) DelAll(prefix string, lock bool) (resp *clientv3.DeleteResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(prefix))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Delete(ctx, prefix, clientv3.WithPrefix())
}

func (es *etcdStore) Lease(key string, ttl int64) (resp *clientv3.LeaseGrantResponse, err error) {
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	if resp, err = es.lease.Grant(ctx, ttl); err != nil {
		return
	}
	return
}

func (es *etcdStore) PutWithLease(key, val string, leaseId clientv3.LeaseID, lock bool) (resp *clientv3.PutResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Put(ctx, key, val, clientv3.WithLease(leaseId))
}

func (es *etcdStore) DelWithLease(key string, leaseId clientv3.LeaseID, lock bool) (resp *clientv3.DeleteResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	if _, err := es.lease.Revoke(context.TODO(), leaseId); err != nil {
		log.Sugar.Error("cancel lease failed. error: ", err)
	}
	return es.kv.Delete(context.TODO(), key)
}

func (es *etcdStore) KeepOnce(leaseId clientv3.LeaseID) (resp *clientv3.LeaseKeepAliveResponse, err error) {
	return es.lease.KeepAliveOnce(context.TODO(), leaseId)
}

func (es *etcdStore) KeepAlive(ctx context.Context, leaseId clientv3.LeaseID) (err error) {
	var (
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)
	// 自动续租
	if keepRespChan, err = es.lease.KeepAlive(ctx, leaseId); err != nil {
		return
	}

	// 处理自动续约应答
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					err = errors.New("stop renewed")
					return
				}
				log.Sugar.Infof("successful renewal,lease id: %d.", leaseId)
			}
		}
	}()
	return
}

func (es *etcdStore) Watch(ctx context.Context, key string, putFunc EventFunc, delFunc EventFunc) {
	var (
		watcherRespChan <-chan clientv3.WatchResponse
		watcherResp     clientv3.WatchResponse
		event           *clientv3.Event
		err             error
	)

	//这里不能用 context.WithTimeout 超时会导致watch退出。
	watcherRespChan = es.watcher.Watch(ctx, key)
	// 处理kv变化事件
	for watcherResp = range watcherRespChan {
		for _, event = range watcherResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				if err = putFunc(event); err != nil {
					log.Sugar.Errorf("watch %s exec put-func failed. err:%s.", key, err)
				}
			case mvccpb.DELETE:
				if err = delFunc(event); err != nil {
					log.Sugar.Errorf("watch %s exec del-func failed. err:%s.", key, err)
				}
			}
		}
	}
}

// 锁
func (es *etcdStore) lock(name string) (func(), error) {
	// 创建一个10s的租约(lease)
	//res, err := eh.client.Grant(context.Background(), 10)
	//if err != nil {
	//	log.Error("lock get grant error. ", err)
	//	return nil, err
	//}
	//ctx,cancel:=context.WithTimeout(context.TODO(),eh.wait*2)
	// 利用上面创建的租约ID创建一个session
	//session, err := concurrency.NewSession(eh.client, concurrency.WithLease(res.name),concurrency.WithContext(ctx))
	session, err := concurrency.NewSession(es.client)
	if err != nil {
		log.Sugar.Error("create lock new session failed. error: ", err)
		return nil, err
	}
	mutex := concurrency.NewMutex(session, name)
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()
	//ctx:=context.Background()

	if err = mutex.Lock(ctx); err != nil {
		mutex.Unlock(ctx)
		session.Close()
		log.Sugar.Errorf("locking key:%s failed. error: ", err)
		return nil, err
	}

	unlock := func() {
		ctx, cancel := context.WithTimeout(context.Background(), es.timeout*2)
		mutex.Unlock(ctx)
		cancel()
		session.Close()
	}
	return unlock, nil
}

func (es *etcdStore) md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

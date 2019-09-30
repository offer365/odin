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
	"sync"
	"time"
)

// etcd 客户端
type etcdStore struct {
	options  *Options
	config   clientv3.Config
	client   *clientv3.Client
	kv       clientv3.KV
	leaseMap sync.Map
	timeout  time.Duration
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
	//es.leaseMap = make(map[string]clientv3.Lease, 0)
	//es.CliM = CliManger{
	//	M: make(map[string]*Cli, 0),
	//	L: sync.RWMutex{},
	//}
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

func (es *etcdStore) Lease(key string, ttl int64) (resp *clientv3.LeaseGrantResponse, err error) {
	lease := clientv3.NewLease(es.client)
	es.leaseMap.Store(key, lease)
	if resp, err = lease.Grant(context.Background(), ttl); err != nil {
		return
	}
	return
}

func (es *etcdStore) PutWithLease(key, val string, id clientv3.LeaseID, lock bool) (resp *clientv3.PutResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	ctx2, _ := context.WithTimeout(context.Background(), es.timeout)
	return es.kv.Put(ctx2, key, val, clientv3.WithLease(id))
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

//func (es *etcdStore) DelAll(prefix string, lock bool) (resp *clientv3.DeleteResponse, err error) {
//	var unlock func()
//	if lock {
//		unlock, err = es.lock(es.md5(prefix))
//		defer unlock()
//	}
//	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
//	return es.kv.Delete(ctx, prefix,clientv3.WithPrefix())
//}

func (es *etcdStore) DelWithLease(key string, leaseId int64, lock bool) (resp *clientv3.DeleteResponse, err error) {
	var unlock func()
	if lock {
		unlock, err = es.lock(es.md5(key))
		defer unlock()
	}
	val, exist := es.leaseMap.Load(key)
	lease, ok := val.(clientv3.Lease)
	if exist && ok {
		if _, err := lease.Revoke(context.TODO(), clientv3.LeaseID(leaseId)); err != nil {
			log.Sugar.Error("cancel lease failed. error: ", err)
		}
	}
	return es.kv.Delete(context.TODO(), key)
}

func (es *etcdStore) KeepOnce(key string, leaseId int64) (resp *clientv3.LeaseKeepAliveResponse, err error) {
	val, exist := es.leaseMap.Load(key)
	lease, ok := val.(clientv3.Lease)
	if !exist || !ok {
		err = errors.New("not found this lease id")
		return
	}
	return lease.KeepAliveOnce(context.TODO(), clientv3.LeaseID(leaseId))
}

func (es *etcdStore) KeepAlive(key, val string) (err error) {
	var (
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
	)
	// 申请一个租约
	lease = clientv3.NewLease(es.client)
	// 申请5秒租约
	if leaseGrantResp, err = lease.Grant(context.Background(), 5); err != nil {
		return
	}
	leaseID := leaseGrantResp.ID
	// 自动续租
	if keepRespChan, err = lease.KeepAlive(context.Background(), leaseID); err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), es.timeout)
	if _, err = es.kv.Put(ctx, key, val, clientv3.WithLease(leaseID)); err != nil {
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
				log.Sugar.Infof("successful renewal,lease id: %d.", leaseID)
			}
		}
	}()
	return
}

func (es *etcdStore) Watch(key string, putFunc EventFunc, delFunc EventFunc) {
	var (
		watcher         clientv3.Watcher
		watcherRespChan <-chan clientv3.WatchResponse
		watcherResp     clientv3.WatchResponse
		event           *clientv3.Event
		err             error
	)
	// 创建一个watcher
	watcher = clientv3.NewWatcher(es.client)
	//这里不能用 context.WithTimeout 超时会导致watch退出。
	watcherRespChan = watcher.Watch(context.Background(), key)
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

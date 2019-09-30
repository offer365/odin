package node

import (
	"context"
	"github.com/offer365/odin/log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

func RunRpcServer(port string, register interface{}) {
	// 注册一个带方法的类型
	if err := rpc.Register(register); err != nil {
		log.Sugar.Error("rpc register failed. error: ", err)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		log.Sugar.Error("net resolve addr failed. error: ", err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Sugar.Error("net listen tcp failed. error: ", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//go rpc.ServeConn(conn) // 并发
		rpc.ServeConn(conn)
	}
}

func GetRemoteNode(ctx context.Context, name, ip, port string, peers []string) (node *Node, err error) {
	var cli *rpc.Client
	dial := func() (ch chan struct{}) {
		ch = make(chan struct{}, 1)
		if cli, err = rpc.Dial("tcp", ip+":"+port); err != nil {
			return
		}
		node = NewNode(name, ip, port, peers)
		if err = cli.Call("Node.Status", Args{name, ip}, node); err != nil {
			return
		}
		ch <- struct{}{}
		if err = cli.Close(); err != nil {
			return
		}
		return
	}
	select {
	case <-ctx.Done():
		log.Sugar.Errorf("call rpc server %s %s:%s timeout. error: %s", name, ip, port, err)
		return
	case <-dial():
		if err != nil {
			log.Sugar.Errorf("call rpc server %s %s:%s failed. error: %s", name, ip, port, err)
			return
		}
		return
	}
}

func GetAllNodes(rpc string,peers []string) (nodes map[string]*Node) {
	var lock sync.Mutex
	var wait sync.WaitGroup
	//var  value atomic.Value
	nodes = make(map[string]*Node, 0)
	//value.Store(nodes)
	wait.Add(len(peers))
	for id, ip := range peers {
		go func(id int, ip string) {
			defer wait.Done()
			name := "odin" + strconv.Itoa(id)
			ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
			n, err := GetRemoteNode(ctx, name, ip, rpc,peers)
			if err != nil {
				log.Sugar.Error("node rpc dial failed. error: ", err)
				return
			}
			lock.Lock()
			nodes[name] = n
			lock.Unlock()
			return
		}(id, ip)
	}
	wait.Wait()
	//sort.Slice(nodes, func(i, j int) bool {
	//	return nodes[i].name < nodes[j].name
	//})
	return
}


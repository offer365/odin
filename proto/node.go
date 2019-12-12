package proto

import (
	"context"
	"encoding/json"
	"sync"

	corec "github.com/offer365/example/grpc/core/client"
	cores "github.com/offer365/example/grpc/core/server"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/log"
	"github.com/offer365/odin/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"time"

	"github.com/zcalusic/sysinfo"
)

var (
	Self          *Node
	auth          *Authentication
	ClientConns   *CliConns
	StaterClients *StaterClis
)

func init() {
	Self = NewNode(config.Cfg.Name, config.Cfg.LocalGRpcAddr())
	ClientConns = new(CliConns)
	StaterClients = new(StaterClis)
	auth = &Authentication{
		User:     _username,
		Password: _password,
	}
}

func NewNode(name, addr string) (n *Node) {
	n = new(Node)
	n.Hardware = new(Hardware)
	n.Hardware.Networks = make([]*Network, 0)
	n.Hardware.Mem = new(Mem)
	n.Hardware.Cpu = new(Cpu)
	n.Hardware.Bios = new(Bios)
	n.Hardware.Host = new(Host)
	n.Hardware.Board = new(Board)
	n.Hardware.Product = new(Product)
	n.Attrs = new(Attrs)
	n.Attrs.Addr = addr
	n.Attrs.Name = name
	n.Attrs.Start = time.Now().Unix()
	n.Hardware.hw()
	return
}

func (hd *Hardware) hw() {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	hd.Host.Machineid = si.Node.MachineID
	hd.Host.Architecture = si.OS.Architecture
	hd.Host.Hypervisor = si.Node.Hypervisor

	hd.Product.Name = si.Product.Name
	hd.Product.Serial = si.Product.Serial
	hd.Product.Vendor = si.Product.Vendor

	hd.Board.Name = si.Board.Name
	hd.Board.Vendor = si.Board.Vendor
	hd.Board.Serial = si.Board.Serial

	hd.Bios.Vendor = si.BIOS.Vendor

	hd.Cpu.Vendor = si.CPU.Vendor
	hd.Cpu.Threads = uint32(si.CPU.Threads)
	hd.Cpu.Model = si.CPU.Model
	hd.Cpu.Cache = uint32(si.CPU.Cache)
	hd.Cpu.Cores = uint32(si.CPU.Cores)
	hd.Cpu.Cpus = uint32(si.CPU.Cpus)
	hd.Cpu.Speed = uint32(si.CPU.Speed)

	hd.Mem.Speed = uint32(si.Memory.Speed)
	hd.Mem.Type = si.Memory.Type

	hd.Networks = make([]*Network, 0)
	for _, val := range si.Network {
		nw := new(Network)
		nw.Speed = uint32(val.Speed)
		nw.Macaddress = val.MACAddress
		nw.Driver = val.Driver
		hd.Networks = append(hd.Networks, nw)
	}

}

func (n *Node) md5() {
	if n.Attrs.Hwmd5 == "" {
		byt, err := json.Marshal(n.Hardware)
		if err != nil {
			return
		}
		n.Attrs.Hwmd5 = utils.Md5sum(byt, nil)
	}
}

func (n *Node) Status(ctx context.Context, args *Args) (*Node, error) {
	n.Hardware.hw()
	n.md5()
	n.Attrs.Now = time.Now().Unix()
	return n, nil
}

func GetAllNodes(ctx context.Context) (nodes map[string]*Node) {
	var lock sync.Mutex
	var wait sync.WaitGroup
	// var  value atomic.Value
	nodes = make(map[string]*Node, 0)
	// value.Store(nodes)
	peers := config.Cfg.AllGRpcAddr()
	wait.Add(len(peers))
	for remoteName, remoteAddr := range peers {
		go func(remoteN string, remoteA string) {
			defer wait.Done()
			if remoteN == Self.Attrs.Name || remoteA == Self.Attrs.Addr {
				Self.Attrs.Now = time.Now().Unix()
				// 重新获取硬件信息
				Self.Hardware.hw()
				Self.md5()
				return
			}
			cli, ok := StaterClients.Get(remoteN)
			if ok && cli != nil {
				n, err := cli.Status(ctx, &Args{Name: Self.Attrs.Name, Addr: Self.Attrs.Addr}, grpc.WaitForReady(true))
				if err != nil {
					if conn, ok := ClientConns.Get(remoteN); conn != nil && ok {
						conn.Close()
						ClientConns.Del(remoteN)
						StaterClients.Del(remoteN)
					}
					NodeGRpcClient(remoteN, remoteA)
					log.Sugar.Errorf("node rpc dial %s %s failed. error: %v", remoteN, remoteA, err)
					return
				}
				log.Sugar.Infof("node rpc dial %s %s success.", remoteN, remoteA)
				lock.Lock()
				nodes[remoteN] = n
				lock.Unlock()
			}
			return
		}(remoteName, remoteAddr)
	}
	wait.Wait()
	nodes[config.Cfg.Name] = Self
	// sort.Slice(nodes, func(i, j int) bool {
	//	return nodes[i].name < nodes[j].name
	// })
	return
}

func AllNodeGRpcClient(peers map[string]string) {
	for name, addr := range peers {
		if name != Self.Attrs.Name && addr != Self.Attrs.Addr {
			NodeGRpcClient(name, addr)
		}
	}
}

func NodeGRpcClient(name, addr string) {
	var (
		Con *grpc.ClientConn
		err error
	)

	Con, err = corec.NewRpcClient(
		corec.WithAddr(addr),
		corec.WithDialOption(grpc.WithPerRPCCredentials(auth)),
		corec.WithServerName(rpcServerName),
		corec.WithCert([]byte(Client_crt)),
		corec.WithKey([]byte(Client_key)),
		corec.WithCa([]byte(Ca_crt)),
	)
	if err != nil {
		log.Sugar.Error(err)
		return
	}

	StaterClients.Put(name, NewStaterClient(Con))
	ClientConns.Put(name, Con)
	return
}

func NodeGRpcServer() (*grpc.Server, error) {
	// Token认证
	auth := func(ctx context.Context) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing credentials")
		}

		var user string
		var pwd string

		if val, ok := md["user"]; ok {
			user = val[0]
		}
		if val, ok := md["password"]; ok {
			pwd = val[0]
		}

		if user != _username || pwd != _password {
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return nil
	}

	// 一元拦截器
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}

	// 实例化grpc Server
	return cores.NewRpcServer(
		cores.WithServerOption(grpc.UnaryInterceptor(interceptor)),
		cores.WithCert([]byte(Server_crt)),
		cores.WithKey([]byte(Server_key)),
		cores.WithCa([]byte(Ca_crt)),
	)
}

// Authentication 自定义认证
// 要实现对每个gRPC方法进行认证，需要实现grpc.PerRPCCredentials接口
type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

type CliConns struct {
	sync.Map
}

func (cc *CliConns) Get(key string) (*grpc.ClientConn, bool) {
	val, exist := cc.Load(key)
	if exist {
		cli, ok := val.(*grpc.ClientConn)
		return cli, ok
	}
	return nil, exist
}

func (cc *CliConns) Put(key string, val *grpc.ClientConn) {
	cc.Store(key, val)
}

func (cc *CliConns) Del(key string) {
	cc.Delete(key)
}

type StaterClis struct {
	sync.Map
}

func (sc *StaterClis) Get(key string) (StaterClient, bool) {
	val, exist := sc.Load(key)
	if exist {
		cli, ok := val.(StaterClient)
		return cli, ok
	}
	return nil, exist
}

func (sc *StaterClis) Put(key string, val StaterClient) {
	sc.Store(key, val)
}

func (sc *StaterClis) Del(key string) {
	sc.Delete(key)
}

package node

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/offer365/odin/log"
	"github.com/zcalusic/sysinfo"
	"net"
	"net/rpc"
	"time"
)

func NewNode(name, ip string) (n *Node) {
	n = new(Node)
	n.Hardware = new(Hardware)
	n.Hardware.Networks = make([]*Network, 0)
	n.Hardware.Mem = new(Mem)
	n.Hardware.Cpu = new(Cpu)
	n.Hardware.Bios = new(Bios)
	//n.Hardware.Chassis = new(Chassis)
	n.Hardware.Host = new(Host)
	n.Hardware.Board = new(Board)
	n.Hardware.Product = new(Product)
	n.Attr = new(Attr)
	n.Attr.Name = name
	n.Attr.IP = ip
	n.Attr.Start = time.Now().Unix()
	return
}

type Node struct {
	*Attr
	*Hardware
}

type Attr struct {
	Name  string `bson:"name" json:"name"`
	IP    string `bson:"ip" json:"ip"`
	Start int64  `bson:"start" json:"start"` // 启动时间
	HwMd5 string `bson:"md5" json:"md5"`
	Now   int64  `bson:"now" json:"now"`
}

type Hardware struct {
	Host    *Host    `bson:"host" json:"host"`
	Product *Product `json:"product"`
	Board   *Board   `json:"board"`
	//Chassis  *Chassis   `json:"chassis"`
	Bios     *Bios      `json:"bios"`
	Cpu      *Cpu       `bson:"cpu" json:"cpu"`
	Mem      *Mem       `json:"mem"`
	Networks []*Network `json:"networks"`
}

type Host struct {
	Machineid    string `json:"machineid"` // 设备id
	Hypervisor   string `json:"hypervisor"`
	Architecture string `json:"architecture"` // 架构
}

type Product struct {
	Name    string `json:"name"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Serial  string `json:"serial"`
}

type Board struct {
	Name     string `json:"name"`
	Vendor   string `json:"vendor"`
	Version  string `json:"version"`
	Serial   string `json:"serial"`
	Assettag string `json:"assettag"`
}

type Chassis struct {
	Type     uint   `json:"type"`
	Vendor   string `json:"vendor"`
	Version  string `json:"version"`
	Serial   string `json:"vSerial"`
	Assettag string `json:"assettag"`
}

type Bios struct {
	Vendor string `json:"vendor"`
}

type Cpu struct {
	Vendor  string `json:"vendor"`
	Model   string `json:"model"`
	Speed   uint   `json:"speed"`
	Cache   uint   `json:"cache"`
	Cpus    uint   `json:"cpus"`
	Cores   uint   `json:"cores"`
	Threads uint   `json:"threads"`
}

// 内存
type Mem struct {
	Type  string `json:"type"`  // type
	Speed uint   `json:"speed"` // 速率
}

type Storage struct {
	Driver string `json:"driver"`
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
	Serial string `json:"serial"`
}

type Network struct {
	Driver     string `json:"driver"`
	Macaddress string `json:"macaddress"`
	Speed      uint   `json:"speed"`
}

func (hd *Hardware) hw() {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	hd.Host.Machineid = si.Node.MachineID
	hd.Host.Architecture = si.OS.Architecture
	hd.Host.Hypervisor = si.Node.Hypervisor

	hd.Product.Name = si.Product.Name
	hd.Product.Serial = si.Product.Serial
	hd.Product.Version = si.Product.Version
	hd.Product.Vendor = si.Product.Vendor

	hd.Board.Name = si.Board.Name
	hd.Board.Vendor = si.Board.Vendor
	hd.Board.Version = si.Board.Version
	hd.Board.Serial = si.Board.Serial
	hd.Board.Assettag = si.Board.AssetTag

	//hd.Chassis.Type = si.Chassis.Type
	//hd.Chassis.Version = si.Chassis.Version
	//hd.Chassis.Vendor = si.Chassis.Vendor
	//hd.Chassis.Serial = si.Chassis.Serial
	//hd.Chassis.Assettag = si.Chassis.AssetTag

	hd.Bios.Vendor = si.BIOS.Vendor

	hd.Cpu.Vendor = si.CPU.Vendor
	hd.Cpu.Threads = si.CPU.Threads
	hd.Cpu.Model = si.CPU.Model
	hd.Cpu.Cache = si.CPU.Cache
	hd.Cpu.Cores = si.CPU.Cores
	hd.Cpu.Cpus = si.CPU.Cpus
	hd.Cpu.Speed = si.CPU.Speed

	hd.Mem.Speed = si.Memory.Speed
	hd.Mem.Type = si.Memory.Type

	hd.Networks = make([]*Network, 0)
	for _, val := range si.Network {
		nw := new(Network)
		nw.Speed = val.Speed
		nw.Macaddress = val.MACAddress
		nw.Driver = val.Driver
		hd.Networks = append(hd.Networks, nw)
	}

}

func (n *Node) md5() {
	if n.Attr.HwMd5 == "" {
		byt, err := json.Marshal(n.Hardware)
		if err != nil {
			return
		}
		h := md5.New()
		h.Write(byt)
		n.Attr.HwMd5 = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
}

type Args struct {
	Name string
	IP   string
}

func (n *Node) Status(args Args, node *Node) (err error) {
	if args.IP == n.Attr.IP && args.Name == n.Attr.Name {
		n.Hardware.hw()
		n.md5()
		*node = *n
		node.Attr.Now = time.Now().Unix()
		return
	}
	return errors.New("ip address error.")
}

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

func GetRemoteNode(ctx context.Context, name, ip, port string) (node *Node, err error) {
	var cli *rpc.Client
	dial := func() (ch chan struct{}) {
		ch = make(chan struct{}, 1)
		if cli, err = rpc.Dial("tcp", ip+":"+port); err != nil {
			return
		}
		node = NewNode(name, ip)
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
		log.Sugar.Errorf("call rpc server %s %s:%s timeout. error: %s", name, ip, port, err.Error())
		return
	case <-dial():
		if err != nil {
			log.Sugar.Errorf("call rpc server %s %s:%s failed. error: %s", name, ip, port, err.Error())
			return
		}
		return
	}
}

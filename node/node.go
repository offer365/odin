package node

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/zcalusic/sysinfo"
	"strconv"
	"time"
)

func NewNode(ip, group, rpc string, peers []string) (n *Node) {
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
	n.Attr.IP = ip
	n.Attr.Group = group
	//n.Attr.Name = name
	for id, ip := range peers {
		if n.Attr.IP == ip {
			n.Attr.Name = n.Group + strconv.Itoa(id)
		}
	}
	n.Attr.Start = time.Now().Unix()
	n.Conf = new(Conf)
	n.Conf.Rpc = rpc
	n.Conf.Peers = peers
	return
}

type Node struct {
	*Conf
	*Attr
	*Hardware
}

type Conf struct {
	Rpc   string   `bson:"rpc" json:"rpc"`
	Peers []string `bson:"peers" json:"peers"`
}

type Attr struct {
	Name  string `bson:"name" json:"name"`
	IP    string `bson:"ip" json:"ip"`
	Group string `bson:"group" json:"group"`
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

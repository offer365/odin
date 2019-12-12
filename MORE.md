## 逻辑
![架构](https://github.com/offer365/odin/blob/master/OdinProcess.png)

[Odin&App交互](https://github.com/offer365/odin/blob/master/Odin&App.png)
## 序列号的生成

> 序列号的结构体
```
// 序列号
type SerialNum struct {
	Sid   string              `json:"sid"`   // 序列号唯一uuid，用来标识序列号，并与 授权码相互校验，一一对应。
	Nodes map[string]*pb.Node `json:"nodes"` // 节点的具体硬件信息。
	Date  int64               `json:"date"`  // 生成 序列号的时间。
}

type Node struct {
	Attrs                *Attrs    `protobuf:"bytes,1,opt,name=attrs,proto3" json:"attrs,omitempty"`
	Hardware             *Hardware `protobuf:"bytes,2,opt,name=hardware,proto3" json:"hardware,omitempty"`
}

type Attrs struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Start                int64    `protobuf:"varint,3,opt,name=start,proto3" json:"start,omitempty"`
	Hwmd5                string   `protobuf:"bytes,4,opt,name=hwmd5,proto3" json:"hwmd5,omitempty"`
	Now                  int64    `protobuf:"varint,5,opt,name=now,proto3" json:"now,omitempty"`
}

type Hardware struct {
	Host                 *Host      `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Product              *Product   `protobuf:"bytes,2,opt,name=product,proto3" json:"product,omitempty"`
	Board                *Board     `protobuf:"bytes,3,opt,name=board,proto3" json:"board,omitempty"`
	Bios                 *Bios      `protobuf:"bytes,4,opt,name=bios,proto3" json:"bios,omitempty"`
	Cpu                  *Cpu       `protobuf:"bytes,5,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Mem                  *Mem       `protobuf:"bytes,6,opt,name=mem,proto3" json:"mem,omitempty"`
	Networks             []*Network `protobuf:"bytes,7,rep,name=networks,proto3" json:"networks,omitempty"`
}
```
> 属性包含 节点名，地址，启动时间，当前时间，硬件md5\
> 硬件信息包含，系统，主板，bios,内存，cpu,mac地址等信息。\
> 将对象序列化成json,转成字节数组，使用 rsa2048+aes256加密生成。\
> ==status 接口其实返回的是整个Node对象==。

## license的生成
> License的结构体
```
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
}

type Attr struct {
	Name  string
	Key   string
	Value int64
}
```
> 解密序列号，检查是否合法（时间），获得sid,硬件md5,这一步是绑定序列号，一个序列号指对应一个授权码。\
> 序列号的合法检查主要是edda的上层应用负责，edda只负责解析序列号，上层应用拿到解析结果后，比对数据库，授权历史里。数据库里是否有一模一样，或部分硬件信息相同的数据，从而分析出客户是否是异常的授权申请。\
> 传入应用信息\
> 将对象序列号成json,转成字节数组，使用 rsa2048+aes256加密生成。

## 激活
> 解密license，检查license的合法性：\
> 节点数量，本机硬件md5是否对应,时间是否合法，是否重复授权，lid与sid 是否一致。

## 重置授权
> 每分钟重置授权核心代码
```
now := time.Now().Unix()
num := (now - lic.Generate) / 60  // 当前的生存周期
// 这里限制了 LifeCycle 只能不断的增大
if num > lic.LifeCycle {
	// atomic.StoreInt64(&(lic.LifeCycle),num)
	lic.LifeCycle = num
} else {
	// atomic.AddInt64(&(lic.LifeCycle),1)
	lic.LifeCycle += 1
}

// 这里限制了 UpdateTime 只能不断的增大  // 授权时间戳
if now > lic.Update {  // bug:？？ now > lic.Update + 59
	// atomic.StoreInt64(&(lic.Update),now)
	lic.Update = now
} else {
	// atomic.AddInt64(&(lic.Update), 60)
	lic.Update += 60
}
```
> 然后将该对象加密保存。\
> Update或LifeCycle 过期，授权将不可用。\
> 存在的问题，比如客户的服务器在授权的时候是准确的，授权成功后，修改成慢时间，应该怎么处理，直接清空授权？\
> 授权一个月，授权成功后，关机或关闭程序，一个月后再运行，如果时间慢了，理论上过期了，但还能使用，这情况无法处理。只能使用 时钟锁

## 加密算法
> 序列号，授权码 rsa2048+aes256加密\
> Auth步骤 请求中Verify字段，默认使用，aes256加密，\
> Cipher使用rsa1024 \
> Auth使用rsa2048 \
> uid使用md5算法 或者sha256(加盐)\
> 以上所有的加密算法都能更换。目前使用多套rsa公私密钥和aes密钥。\
> odin与edda 的部分常量，用户密码，哈希盐，加密算法，在对接中要约定好。

## 补充上次分享
> 配置功能\
> 建议使用gRpc，安全性和效率上要优于Restful。

#### etcd
> etcd是一个golang编写的==分布式==、==高可用==的==强一致性键值存储==系统，用于提供可靠的分布式键值(key-value)存储、配置共享和服务发现等功能。etcd可以用于存储==关键数据==和实现分布式调度。\
> etcd基于raft协议，通过复制日志文件的方式来保证数据的强一致性。在etcd之前，常用的是基于paxos协议的zookeeper。


> 一：服务发现\
所有服务将元信息存储到以某个 prefix 开头的 key 中，然后消费者从这些 key 中获取服务信息并调用。消费者也可以 watch 这些 key 的变更，以便在服务增加和减少时及时获得通知。\
> 二：配置共享\
应用将配置信息存放到 etcd，当配置信息被更改时可以通过 watch 机制从 etcd 及时获得通知。\
> 三：分布式锁\
由于 etcd 中的数据是一致的，当多个应用同时去创建一个 key 时只有一个会成功，创建成功的应用即获取了锁。\
> 


> ==Kubernetes==就是使用了etcd存储持久化数据。\
> ==TiDB==(开源分布式数据库)一致性实现也是使用etcd的raft实现。\
> 在==微服务==中常常被用于服务发现。\
> ==CoreOS==(内置Docker容器的操作系统/是为云而生的操作系统)etcd 以默认的形式安装于每个 CoreOS 系统之中。Etcd是CoreOS生态系统中处于连接各个节点通信和支撑集群服务协同运作的核心地位的模块。\
> ..............


#### 加密狗
[时钟锁](https://item.jd.com/22221174281.html)\
[max超级锁](https://item.jd.com/22057174395.html)\
![image](https://img10.360buyimg.com/imgzone/jfs/t1/25873/21/5508/363702/5c3ee1f5E4070517b/910a47b55795cfb6.jpg)

## 编译参数

```
# 对外发布的程序
go build -ldflags '-s -w'
```

- 压缩程序
- 禁止gdb调试

## 基于软件的保护方式：
- 注册码
- 许可证文件
- 许可证服务器
- 应用程序服务器

主要破解手段及应对措施：
> 1、编译型语言：C/C++,Delphi\
> 破解手段：\
> 使用反汇编工具，静态分析汇编代码\
> 使用调试工具动态跟踪，找到加密点，暴力破解。\
> 制作文件补丁或者内存补丁。
> 
> 应对措施：\
> 既调用 API，又加壳保护\
> 隐藏加密点\
> 避免简单的 License 检查\
> 使用加密锁的查询或加密解密功能\
> 检查模块的完整性，可以考虑使用数字证书对所有可执行文件进行签名，在
> 程序运行过程中验证签名的正确性\
> 如果是 C/S 或 B/S 架构的应用程序，可以将部分重要功能或检查放在服务器端\
> 综合运用多种技术手段

- 防动态调试
- 防内存DUMP
- 防动态注入

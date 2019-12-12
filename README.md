# odin #

----

## what this?
- 应用在局域网中或互联网中的分布式授权服务。
- odin 是一个 license server  用于给多种应用的多个客户端提供授权服务。
- 目前只能运行在 linux 系统中。
- 这个应用需要与[edda](https://github.com/offer365/edda) 配合使用。
- 在 [edda](https://github.com/offer365/edda) 中生成授权码，在odin中激活。多个客户端或应用从odin 获取授权信息。
- 被授权的应用可以通过 restful 或 gRPC 与 odin 交互。
- 更多内容请阅读 [MORE.md](https://github.com/offer365/odin/blob/master/MORE.md)



## 特点
- 分布式。
- 可以通过绑定应用的硬件信息鉴权。
- 支持 https 加密通信。
- 部署简单，无依赖。

#### odin
> odin 可以根据 客户端提供的token 限制产品安装的位置。

#### edda
> edda 不提供授权管理，审批流程，客户关系等功能\
> 只实现了授权相关四个核心功能。

**gRpc接口**

```
service Authorization {
    rpc Resolved (Cipher) returns (SerialNum); // 解析序列号
    rpc Authorized (AuthReq) returns (AuthResp); // 授权
    rpc Untied (UntiedReq) returns (Cipher); // 解绑
    rpc Cleared (Cipher) returns (Clear);  // 清除
}
```


## 技术栈
- 使用 GoLang 语言开发。
- 使用 gin web 框架。
- 嵌入etcd。
- 支持 RestFul 和 GRPC 传输数据。
- 前后端分离 bootstrap + jquery + ajax 。

## 安全
- RestFul 和 gRPC 都使用 https 传输数据。
- 序列号和授权码使用 aes32+rsa2048加密。Auth步骤可以采用多种加密方式。
- 通过 ==到期时间== 与 ==生存周期== 相互印证防止时间篡改。
- odin 不会存储明文，存储加密过或hash后的值。
- N个节点(奇数)部署时，可以允许 (N-1)/2 个节点挂掉，且数据不会丢失。
- 内嵌 etcd 支持 用户认证。

## 性能
- CPU:Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz * 2
- MEM: 2G

> 单节点下1000个客户端同时访问。cpu:27%  mem:5%\
> 三节点下1000个客户端同时访问。cpu:33%  mem:7%



## 应用交互流程
在demo/proto/auth.proto 包里定义了request 和 reponse 消息体

###### request
```
message Request {
    string app = 1; // 应用名称
    string id = 2; // 应用ID
    int64 date = 3; // 当前时间戳,客户端服务端误差不能超过600s
    string verify = 4; // 该字段是密文,用于校验request的参数,客户端要根据 app,id,date,token(唯一且固定的值)加密生成;eg: {"app":"nlp","date":1571987046,"id":"app01","token":"xxxxxx"} 对此加密
    string umd5 = 5; // response中返回的 data.cipher 解密后的值的md5;在active 步骤中此参数无效,仅在 keepline 与 offline中有效
    int64 lease = 6; // 租约ID 在active 步骤中此参数无效,仅在 keepline 与 offline中有效
}
```
###### response
```
message Data {
    bytes auth = 1; // 解密后是应用的一些属性;eg:{"attrs":[{"Name":"热词","Key":"hotword","Value":111},{"Name":"类热词","Key":"classword","Value":111}],"time":1571994906931717352}
    int64 lease = 2; // 租约ID
    bytes cipher = 3; // 加密UUID生成的密文。
}

message Response {
    int32 code = 1; // 返回状态码 200 OK;
    Data data = 2;
    string msg = 3; // 错误 或 成功的消息
}
```

#### Auth 认证
**request**
> request 必须包含 app,id,date,verify参数。\
> verify: 该字段是密文,用于校验request的参数,app端要根据 app,id,date,token加密并base64编码生成。\
> eg: {"app":"nlp","date":1571987046,"id":"app01","token":"xxxxxx"}\
> token要求唯一且不可变，app端可以根据硬件信息，比如:Mac地址，机器码，主板编号，硬盘编号等生成。\
> odin不会记录token的值，而是保存token的hash值，用来校验app端的身份。\
> 一个app端只能有一个token,如果token变了需要解绑。\
> odin会解密并校验verify的值与时间戳，值必须相等，时间戳误差不能
超过600s.

**response**
> 通过验证后，odin会将UUID和应用属性加密返回给app,分别是 response.data.cipher 和 response.data.auth字段\
>例： 
```
{"attrs":[{"Name":"热词","Key":"hotword","Value":111},{"Name":"类热词","Key":"classword","Value":111}],"time":1571994906931717352}
```
> data 中包含 lease 表示租约id,这个字段会在keepline中用到。\
> response.data.cipher 和 response.data.auth 建议采用不同的加密方式。\
> 并且约定好 hash算法 和 是否加盐。
> app端收到 response 后解密，得到 uuid和认证信息 

```
sequenceDiagram
App->>Odin: request消息
Odin->>App: response消息
```

#### 在线 Keepline
**request**
> app端将 在auth中获得的uuid hash计算后与 lease 一起发送给odin，分别对应request.umd5和reques.lease 字段。requset 中verify字段留空即可。其他的与auth步骤中一样。
>

**response**
> odin会校验umd5值与lease，然后续租。如果app端超过10秒未发送数据。会认为app端下线，将app的信息删除。此时app端需要重新进行auth步骤。
> odin仅返回lease id。

```
sequenceDiagram
App->>Odin: request消息
```


#### 下线 Offline
**request**
> app端将 在auth步骤中获得的uuid hash 后与 lease 一起发送给odin，分别对应request.umd5和reques.lease 字段。requset 中verify字段留空即可。其他的与auth步骤中一样。
>

**response**
> odin会校验umd5值与lease，然后将app的信息删除，认为app端下线。\
> odin仅返回lease id。

```
sequenceDiagram
App->>Odin: request消息
```

#### 绑定与解绑
**绑定**
> 即：将 app/id : token  的对应关系存储到硬盘。\
> 当一个新的app端发起auth认证时，如果是第一次认证就注册该app的id与token。\
> 如果不是第一次就检查该app的id与存储的token是否对应，直到该应用的实例用尽。\
> 所以要求app端生成的token固定且唯一。\
> 该token可以是根据某个硬件信息生成，可以是密文或者明文。\
> 总之odin 并不关心该token 的值，存储的仅仅是token的hash值。\
> 对于不想使用绑定功能的应用可以 将token的值设为 "app/id" eg: "nlp/app01"。\
> token 需要放在 auth 步骤的verify字段中。\
> 在实际存储中app/id与toekn的对应关系可能是  
```
eg: /NcJEs2UCgYEA/mh0SYJIGyadW/hash(app,salt)/hash(id,salt)    hash(token,salt)
```

**解绑**
> 解绑常常发生在app端硬件信息改变或app端迁移或其他问题导致token值发生改变。\
> 需要将 app/id 与token 对应关系解除。\
> 在 [edda](https://github.com/offer365/edda) 中输入 app/id,会产生一段密文。\
> 这个密文里包含该app/id 在odin在的路径，odin解密后，比对路径与时间戳。\
> 如果合法就将该 路径 与对应的token删除，从而解绑。app端再发次auth认证即可注册新的token。

## 安装运行 ##
#### 安装odin

```
unzip odin-xxx-linux.amd64.zip
cd odin
sh install.sh
# 请先修改odin.yaml,appctl.sh 中的IP地址为本机IP。
./appctl.sh resetcode 
./appctl.sh getcode
```

>
> 访问 https://127.0.0.1:9527


#### 相关说明
> 配置文件是 odin.yaml。\
> 修改odin.service 可以指定程序与配置文件的位置。\
> appctl.sh 封装了linux下的一些常用api操作。\
> app端的测试，请参考 demo。\
> 建议使用三个节点提供服务。

## Demo
- 目前仅实现了go语言的restful 和 grpc demo。在demo 目录下。

## TODO
* 序列号/授权码采用 椭圆曲线ECC加密算法 (maybe)
* 支持 圣天诺Sentinel Time时钟锁 (maybe)
* 前端交互支持 session，不再使用 BasicAtuh。Linux下，curl访问时 根据 UserAgent区分 使用 BasicAtuh 认证。(maybe)
* 提高测试代码覆盖率。
* 内嵌 etcd 支持https。
* 支持 windows。实现获取硬件信息接口。(maybe)

## 使用介绍 ##
1. 先安装odin 并运行。访问web端口，默认账号密码：admin:123 可在配置文件odin.yaml 中修改。
2. 使用web 或访问api 接口生成 序列号。
3. 在 [edda](https://github.com/offer365/edda) 里根据约定新建应用，并配置该应用的属性。
4. 将序列号在 [edda](https://github.com/offer365/edda) 解析，根据实际情况填写，生成license。
5. 在 odin 的激活页面或相关api 导入序列号激活。
6. 客户端或应用访问 odin 的client api。详细示例见API.md和demo


## License
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

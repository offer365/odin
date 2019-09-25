# odin #

----

## what this? ##
- 应用在局域网中或互联网中的分布式授权服务。
- odin 是一个 license server  用于给多种应用的多个客户端提供授权服务。
- 通过内嵌 etcd 来存储数据。
- 目前只能运行在 linux 系统中。
- 这个应用需要与[edda](https://github.com/offer365/edda) 配合使用。
- 在 [edda](https://github.com/offer365/edda) 中生成授权码，在odin中激活。多个客户端或应用从odin 获取授权信息。
- 其他应用通过访问 odin 的api 接口来获取信息。

## 安装运行 ##
#### 安装odin

```
unzip odin-xxx-linux.amd64.zip

cd odin

sh install.sh

./appctl.sh resetcode 

./appctl.sh getcode

```

>
> 访问 https://127.0.0.1:8888


#### 相关说明
> 配置文件是 odin.json 
>
> 修改odin.service 可以指定程序与配置文件的位置
>
> appctl.sh 封装了linux下的一些常用api操作。
>
> 客户端的测试，请参考 client_example.
>
> 建议使用三个节点提供服务。



## 使用介绍 ##
1. 先安装odin 并运行。访问web端口，默认账号密码：admin:123 可在配置文件odin.json 中修改。
2. 使用web 或访问api 接口生成 序列号。
3. 在 [edda](https://github.com/offer365/edda) 里根据约定新建应用，并配置该应用的属性。
4. 将序列号在 [edda](https://github.com/offer365/edda) 解析，根据实际情况填写，生成license。新的license 可在授权历史中找到。直接复制
5. 在 odin 的激活页面或相关api 导入序列号激活。
6. 客户端或应用访问 odin 的client api。详细实例见API.md 
7. 客户端解密auth字段的值,获取相关的授权信息。解密cipher字段的值来获取该应用的uuid。
8. 客户端携带uuid与租约id发送心跳请求，10秒不发送，会删除客户端实例。需要重新访问client api。
9. 客户端退出时可发送注销请求。

   
   
## License
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).





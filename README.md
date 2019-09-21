# odin #

----

## what this? ##
- 应用在局域网中或互联网中的分布式授权服务。
- odin 是一个 license server  用于给多种应用的多个客户端提供授权服务。
- 通过内嵌 etcd 来存储数据。
- 目前只能运行在 linux 系统中。
- 这个应用需要与[edda](https://github.com/offer365/edda) 配合使用。
- 在 [edda](https://github.com/offer365/edda) 中生成授权码，在odin中激活。
- 其他应用通过访问 odin 的api 接口来获取信息。

## 使用手册 ##
1. 安装odin
> `cd /home/admin`
>
> `go get github.com/offer365/odin`
>
> `cd odin;go build`
>
>`cp scripts/odin.service /usr/lib/systemd/system/`
>
>`systemctl enable odin`
>
>`systemctl start odin`
>
> 访问 ip:8888


2. 相关说明
> 配置文件是 odin.json 
>
> 修改odin.service 可以指定程序与配置文件的位置
>
> appctl.sh 封装了linux下的一些常用api操作。
>

## License
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).





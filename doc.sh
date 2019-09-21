# token
token1  = `f65490272577bbf04ef029ea7a4669230b`
# 认证
urlAuth = `http://10.0.0.200:9999/odin/api/client/collection/auth`
# 心跳
urlHeartbeat   = `http://10.0.0.200:9999/odin/api/client/collection/heartbeat`
# 关闭
urlShutdown   = `http://10.0.0.200:9999/odin/api/client/collection/shutdown`

访问流程
1, 身份验证
客户端 携带id与token    发送认证请求  到服务器 
例: http://10.0.0.200:9999/odin/api/client/collection/auth?id=bjszy001&token=f65490272577bbf04ef029ea7a4669230b
服务器正常返回:
{"code":200,"conf":"ZVjJ+eA/qXBpmSQuMEl2d6rLpCt8X9JvMo8G9fNK3nBX4pXFM6vYUI69yfk6fsqj2VT6Mg1GJhY1Bch+aFGG/4u+/iKo+re69LRJyYSZyaOap0mDlmE3B6upSeBt63Zi/R+OnHBd4uCncT+wVZ//6Bu0WHnjm7iCgQzdNoXrk145/9HEwdlXM/8xEZevS5N/Zxht2P0OjOE86rljintl6cKCNL+KGt4o3LpKKBZzUte8GxHBPl68i2uK6hVzeSVNYk/PmOINt+LgMk+nnJKgIk+hAjXwBwzOl9yCf7ZWgedY3rCjGEoyJ7UzlC2RUFPAU327V9iI/91+CVaYOq+HNqSueZURon/zFatwESK8vd4Sjh+7f8q6UK//m6ba5KdE2MfBxQhDuAyy6OjCZXuuF9XAQCb6EgIjHQqtjzD9CN+VlUY1cXJxsbPCj/S7cCMNegZIi9nRM+lZ9hMPcmqgLjBOfEQ351hdifDkXQd4NmNOW5FLm9fWgiAsOm6TTmzSvzpapczJ8wgoUVGqhBJmV4hHA+thTvVP09+sX9ev3J/GaQsMSU5c/pUwhNzgYwK8DDvExZTKv6ssfPlbeXvsVuXOfRhfkELfV5Fr3rYQsBQboLuvPbZPbtW+nbTmcaKqS0RvlB+8Ne+vtnNlVMMxVg+wS/LmKiMtC1GmSIDT3z+fbCyuq2KhVQbhm0yJ8H4Y","expire":"2019-03-25 14:17:33","lease":1686714280076933163,"msg":"auth: Verification passed."}
code 200 表示请求成功，lease 表示 该客户端实例的租约id.  msg与err 分别表示正确与错误的信息。conf 是加密的配置信息，需要解密。
conf 要经过 aescbc 解密，与rsa 公钥解密。
解密后：
{"serverpath":"http://10.10.10.166:8000","audiopath":"10.10.10.166","websocketpath":"ws://10.10.10.166:8080","studypath":"10.10.10.166","asrpath":"10.10.10.166","upload":"1","token":"1234","place":"testtest","channels":[5,40],"c":"1","record":"","identity":"ba131a92-5860-43b9-ac36-383eb29e847a"}
解密后有个字段 identity 表示该客户端本次请求的uuid

服务器异常返回:
{"code":400,"expire":"0","lease":0,"msg":"","err":"auth: The maximum authorized instance has been reached."}  超过授权个数
{"code":400,"expire":"0","lease":0,"msg":"","err":"auth: This instance already exists."}  该客户端实例已经存在
{"code":500,"expire":"0","lease":0,"msg":"","err":"auth: Failed to get the configuration of the modified client."}  获取客户端配置失败
{"code":500,"expire":"0","lease":0,"msg":"","err":"auth: auth: Authorization encryption failed."}  服务器加密配置文件失败
{"code":500,"expire":"0","lease":0,"msg":"","err":"auth: auth: The client instance failed to put."}  服务器put客户端失败
{"code":400,"expire":"0","lease":0,"msg":"","err":"auth: Token verification failed."}  验证token 失败。
{"code":404,"expire":"0","lease":0,"msg":"","err":"default: bad Request."} url 错误 
{"code":400,"err":"Authorization expired.","expire":"0","lease":0,"msg":""}  授权到期
2, 存活心跳
客户端 携带 id lease identity 发送 请求  到服务器
例: http://10.0.0.200:9999/odin/api/client/collection/heartbeat?id=bjszy6&identity=32be8995-e3f7-4735-ace2-2396b90df529&lease=1686714280076932928
该请求需要每3秒发送一次，如果服务器10秒没有收到请求，就认为该实例挂了。需要重新发起 认证请求。
服务器正常返回:
{"code":200,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"heartbeat: Successful renewal","err":""}
服务器异常返回:
{"code":400,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"heartbeat: Lease id error.","err":""} 租约id错误
{"code":400,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"heartbeat: Renewal failure.","err":""}  续租失败
{"code":400,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"heartbeat: Failed verification, or unknown identity.","err":""} 验证失败或身份未知。
{"code":404,"expire":"0","lease":0,"msg":"","err":"default: bad Request."} url 错误 Authorization expired.
{"code":400,"expire":"0","lease":0,"msg":"","err":"Authorization expired."}  授权到期

3,客户端关闭
客户端 携带 id lease identity 发送 请求  到服务器
例: http://10.0.0.200:9999/odin/api/client/collection/shutdown?id=bjszy7&identity=55d73833-e23a-45e2-98b5-ebdc3cae9c59&lease=1686714280466675843
服务器正常返回:
{"code":200,"expire":"2019-03-25 14:42:00","lease":1686714280466675855,"msg":"shutdown: Deleting an instance succeed."}
服务器异常返回:
{"code":400,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"shutdown: Lease id error.","err":""} 租约id错误
{"code":500,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"shutdown: Deleting an instance failed.","err":""} 服务端删除实例失败
{"code":400,"expire":"2019-03-25 14:16:23","lease":1686714280076932928,"msg":"shutdown: Failed verification, or unknown identity.","err":""} 验证失败或身份未知。
{"code":404,"expire":"0","lease":0,"msg":"","err":"default: bad Request."} url 错误 Authorization expired.
{"code":400,"expire":"0","lease":0,"msg":"","err":"Authorization expired."}  授权到期

###
客户端最简配置文件
{"id":"bjszy001","server":["10.0.0.200:9999","10.0.0.201:9999","10.0.0.202:9999"]}
客户端应该均衡的访问其中一台服务器（随机，或者散列计算），一个完整的访问流程，建议在一台服务器上完成。也支持在不同的服务器上进行。
客户端应该有断开重连的机制。断开指 心跳超时后，再次进行auth 请求。
分布式集群中，三个节点服务器最多容忍一个节点宕机。
授权服务具备配置中心功能，可以根据不同的客户端id,返回不同的配置信息。
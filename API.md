# Odin Api 文档


# *序列号*
### 获取序列号
`curl-k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/server/code`
###### 返回示例:
`{
     "code": 200,
     "data": {
         "date": "2019-08-28 14:53:33",
         "key": "wjUtPchpl2oRLPeno/57CNIg10KN169CAm9Q9vriJhX1V+Aco90/Wpt+1e6uhsv7cxOcOFj4vSE60tQSZE6keszB+Jrx795XTuOoodC71z5wrtKCp6TL5yEFKQ+h9E2YRixcS+3zUDV+y4sOt0cGTVdYeCSUKBVS+T4mNhUjbTxuGNQkMGA+bNH4WtNiJM2sE+uClt7vfkq7c2ThKEr/tNa6HS7E3G4sdLTlPTSaRWBo8yt2JPsO4iy8sCqJExQo/29/j13j5wKc45/FJwuLJuYG1nE+oYHqXjE4LqUD/HzUUVSz/1tdCBfg7uwT9EJScedSFdu0LR0X8qmRekP8O4ywX3HkWTs8/wrWEMyYU7YgHfevAzkdt7HhrohzFTbDCPj6J78hg+KiK6iSa5CGiPSQVyGSVpfWijFMyXf3wH2yc/m23+nrSG0eqlmBhHdR3cQ9QwjiHD9CEIB55ImC3zBtirr3TE9moihVWqgzxgsKWbEnWPFEW5VHDStGLJooC+kl+0wC4Y/yZLUnyw1XEUOlBhG/9lDkrh6Nhdt9bRz0Xpql+b+Rox6pWuML/D3tCnKxNB530zpNss6AJ9jDvECO2S8l5M3M+iVFrpjN2LVFIjKdtfqjCWsCibNJhIHlsi9WeBv8QQzmlCCPF8ou+/ycvXDEfHOlozsVaWEvC2zhTv6EZG82GInss31ff7GV"
     },
     "msg": "在授权成功前请勿重启进程或系统，否则序列号将变更。请保证机器硬件和系统时间正确(误差5分钟内)，否则可能会导致进程异常或者授权失效。"
 }`

### 重置序列号
`curl -k -s -X POST --user admin:123 https://10.0.0.200:8888/odin/api/v1/server/code`
###### 返回示例:
`{
     "code": 200,
     "data": {
         "date": "2019-08-28 14:54:25",
         "key": "ke7bWZnNDOP5BhsRRG2PyBKzFmZb0YkpFh476HqeWKGAp8sb3N6r34hnV8RNs+p0l1mDzCAL4wZd9CoCwCcusXl6AQy6cMBELnmAZC3HNLgECTJLaK5L7SOYIsKTCodvHMedwHAf+rXkoEP5qGVn9fuFTe/5sK1TnPqZZcLPhpf9OoxPR+xuq17/j6iCjvSdu9xzWR3vw9oJ4jjdIEvHV08PmBZiN2mcjGIcKVQ6Z9Zhc/wqQERQm/SdI4oGvQQ6FlrneWtCE+wx2TbT3lJqatfIbwCiLJIUJpOIwffJkBFnRa7ob3gErVqBcwdDAcJpDfaLVbtQl9j5tbgsng7Zyvj1xIf84wfebKzYWMHKl/pq3Sa0w6Mz/150WIK+Z7OCu4pE/27CC2zNY5EWxHcrEy5ugUzlyVpSMBk/kwOeJdTp9uTzg3u6F4IVXu0Njw53WG4F+qIA2R49SblytzCCM6swS9f8Bz6Kv6QnZ48iWmiQvO279XImdYoU6nyfzwZEi9OC4qAuZ2KWL0FWIsQyf7E2B1lIC0jU130P4UvaWDbTunGJnMCW8RNXJ4ElMG6BbXIS3uJhcDbMG+Lj5Vr+HzOQXV9a1Gaj518BUoRJp2fvMD2yg7DYV0UhVnkA9fPULUbjumoeXImBVYc2IrZOgN2fKqCf++PogH4eaZ19ipdRHx2Mtuc6W/d+AhWDVfGX"
     },
     "msg": "在授权成功前请勿重启进程或系统，否则序列号将变更。请保证机器硬件和系统时间正确(误差5分钟内)，否则可能会导致进程异常或者授权失效。"
 }`

### 获取序列号二维码
`curl -k -s -X GET --user admin:123 -o qr-code.jpg https://10.0.0.200:8888/odin/api/v1/server/qr-code`
###### 返回示例:
`二维码图片`

---

# *授权码*
### 导入授权码
`curl -k -s -X POST --user admin:123 -F key="授权码" https://10.0.0.200:8888/odin/api/v1/server/license`
###### 返回示例:
` {"code":200,"msg":"激活成功。"}`

### 查看授权信息
`curl -k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/server/license`
###### 返回示例:
`{
     "code": 200,
     "data": {
         "update_time": "2019-09-19 17:22:00",
         "life_cycle": 15,
         "apps": [
             {
                 "title": "应用1",
                 "data": [
                     {
                         "title": "class-word",
                         "value": "11"
                     },
                     {
                         "title": "hotword",
                         "value": "11"
                     },
                     {
                         "title": "model",
                         "value": "11"
                     },
                     {
                         "title": "最大实例",
                         "value": "11"
                     },
                     {
                         "title": "到期时间",
                         "value": "2019-10-03 00:00:00"
                     },
                     {
                         "title": "最大生存周期",
                         "value": "19612"
                     }
                 ]
             },
             {
                 "title": "测试应用",
                 "data": [
                     {
                         "title": "bingfa",
                         "value": "11"
                     },
                     {
                         "title": "connect",
                         "value": "11"
                     },
                     {
                         "title": "最大实例",
                         "value": "11"
                     },
                     {
                         "title": "到期时间",
                         "value": "2019-10-03 00:00:00"
                     },
                     {
                         "title": "最大生存周期",
                         "value": "19612"
                     }
                 ]
             }
         ]
     },
     "msg": "success"
 }`

### 注销授权
`curl -k -s -X DELETE --user admin:123 https://10.0.0.200:8888/odin/api/v1/server/license`
###### 返回示例:
`{"code":200,"msg":"mVxvD8OmBo8Jjco2UDz+BEqM68H1dBcz77hB/g61tBkEJ+fsDxsNkWCC/mgitVoDs01OS3y9QYgNFTBHPk2NFxNHSB1vQkB5awvjW6oKxwBo8Hq2ISyp+X9feIt5nX+jJwEqFenGJB4fFMWrYDJDE3DkZ19WnDRpu9av03n/MFKBtwZAvvi8IdJ7PQcMw1AzK98zg9Y7rY3K0Sd18UTvmO6J5ZpUp6qzpES3Q1KSRF332AV9Wl0wEd68WS8y+pIMbQ+Z3pb97vjWRbagsps8+8K5mDaS0j6hYQP5dqDkkvbMfsM1UkCbfhu65D0rh7Z1Ok2dSTp9Ps/D1xvrB2fGNvu1kQ9WYBnE8LTvSpSv6haNdysn+uxVaiyjyVo6n2YpyaPd/ZrJcX0KzNFHLYwwAQqCa59udiMpX0TdA/GrUYSc+n+5vaywLHQ28A0kJfgBGttttyYV2nglvybxwIIZkpfXp6pLmf50vINaVT+dwX+QD24cePt/KKBHZPM0dAPwPoGf2Wh8VPHFcDpMSpaeuBPy3wg6cnfl+NdZWwNF6UpKNy6PdJADOKowPZhJZnxXi59pMeuGEV5akW2eF9KmnFKuyBU1ZTte58/ttZ3T46xDE3n1QZJgN2JmeHO4cTKnnto/1SUtsk2f9HyKhonP0oYxdCtNOepSL9aqtw7t4rbQFbLjRbbfGc2vNzX77hmA"}`

### 获取注销二维码
`curl -k -s -X GET --user admin:123 -o qr-license.jpg https://10.0.0.200:8888/odin/api/v1/server/qr-license`
###### 返回示例:
`二维码图片`

---

# *节点状态*
### 查看节点状态
`curl -k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/server/nodes`
###### 返回示例:
`{
     "code": 200,
     "data": [
         {
             "id": "odin0",
             "online": "节点:odin0 ip:10.0.0.在线 00天00小时04分钟09秒."
         },
         {
             "id": "odin2",
             "online": "节点:odin2 ip:10.0.0.202 在线 00天00小时17分钟54秒."
         },
         {
             "id": "odin1",
             "online": "节点:odin1 ip:10.0.0.201 在线 00天00小时17分钟57秒."
         }
     ],
     "msg": "success"
 }`

---

# *配置接口*
### 新增配置
`curl -k -s -X POST --user admin:123 -F text=demo https://10.0.0.200:8888/odin/api/v1/client/conf/{name}`
###### 返回示例:
`{"code":200,"msg":"Post or Put key success."}`
### 删除配置
`curl -k -s -X DELETE --user admin:123 https://10.0.0.200:8888/odin/api/v1/client/conf/aa`
###### 返回示例:
`{"code":200,"msg":"Delete key success."}`
### 修改配置
`curl -k -s -X PUT --user admin:123 -F text=foobar https://10.0.0.200:8888/odin/api/v1/client/conf/aa`
###### 返回示例:
`{"code":200,"msg":"Post or Put key success."}`
### 获取配置
`curl -k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/client/conf/aa`
###### 返回示例:
`{"code":200,"data":[{"name":"aa","text":"bb"}]}`
### 获取所有配置
`curl -k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/client/conf/`
###### 返回示例:
`{"code":200,"data":[{"name":"aa","text":"foobar"},{"name":"bb","text":"demo"}]}`

---

# *Client接口*
### 获取认证
`curl -k -s -X POST --user admin:123 https://10.0.0.200:8888/odin/api/v1/client/auth/demo1/aa`
###### 返回示例:
`{
     "code": 200,
     "data": {
         "auth": "CnS6v03p35u27XbfCCF4WdN+h9dkyk6NUrC4dXPPgAh2nWrcfGUV9DTEpSLlr5OzhZj6I23bfBfXCTgIrMWNJg3cH+P2PlpqfX8rNacwtUsZx35toFYULm021RtVwI/z56bID7xZUh+FEb1in3AgbuHydy7UZeuyYV/vi4px5IuJTKIL099tbakHE/OPYNExw+1qaM5tlC7rqbDoiFO1q9tb4WEM9eWJn7UU/wll2ooQFtCPcPR4nmaNkavzsXSqfK1V5+GGX6o5RiKaU6j7FK2oNVCLUpJQ4+TMSZ25JxfkmcbKtV6LQaMz8AokF7oDbgTNskvTdsiTIzUtoR+7WSR3hl7USerPDytHeIxpgl8=",
         "lease": 6479674172846204226,
         "cipher": "FOrMMgjfKSpk7jvIirnuXK7uPdey0H3PpUoV1aIQcFBdPlHXa1VJZGhU4VH/oqrfCBbqMcOP9yohYK64Rd/mfWhD5IMTtYN+fMRsr0H+r0p28nkNrxoB6j1qwsoGOAn7/JEE643yBq3ZKlWXBVnyu8JXT3VIGqId3zWihTtA7pZaZi8cFZkx6msv3jZ6paUYbYOjuIdvpcfwodRSQbTMNFduAilh09epwOieW8z6JbVTH3YjdeJBjvZ9ReN/wH1HIPFJ9e1/LlfpTU5AC7fb06v+pXqNV3UcFcDRZcHcboJcyKCs2i5azrM6U9xVT06dsyIFZUyeY00ikarXUQKp8qpNvTPI+jDKtCwtiYg4qDg="
     },
     "msg": "success"
 }`
### 心跳
`curl -k -s -X PUT --user admin:123 -d '{"uid":"741f9919-27e6-49d0-a6f8-1116a713f271","lease":6479673129519463168,"auth":""}' https://10.0.0.200:8888/odin/api/v1/client/auth/app1/aa`
###### 返回示例:
`{"code":200,"data":{"auth":"","lease":6479674172846204816,"cipher":""},"msg":"Successful renewal."}`
### 关闭
`curl -k -s -X DELETE --user admin:123 -d '{"uid":"741f9919-27e6-49d0-a6f8-1116a713f271","lease":6479673129519463168,"auth":""}' https://10.0.0.200:8888/odin/api/v1/client/auth/app1/aa`
###### 返回示例:
`{"code":200,"data":{"auth":"","lease":6479674172846204672,"cipher":""},"msg":"Deleting an instance succeed."}`

---

# *Client在线信息接口*
### 在线信息
`curl -k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/client/online/demo1`
`curl -k -s -X GET --user admin:123 https://10.0.0.200:8888/odin/api/v1/client/online/`
###### 返回示例:
`{
     "code": 200,
     "data": [
         {
             "id": "app1/a0",
             "info": "节点:a0(10.0.0.254) app1 在线 00天00小时00分钟05秒."
         },
         {
             "id": "app1/a1",
             "info": "节点:a1(10.0.0.254) app1 在线 00天00小时00分钟05秒."
         },
         {
             "id": "app1/a10",
             "info": "节点:a10(10.0.0.254) app1 在线 00天00小时00分钟00秒."
         },
         {
             "id": "app1/a2",
             "info": "节点:a2(10.0.0.254) app1 在线 00天00小时00分钟04秒."
         },
         {
             "id": "app1/a3",
             "info": "节点:a3(10.0.0.254) app1 在线 00天00小时00分钟04秒."
         },
         {
             "id": "app1/a4",
             "info": "节点:a4(10.0.0.254) app1 在线 00天00小时00分钟03秒."
         },
         {
             "id": "app1/a5",
             "info": "节点:a5(10.0.0.254) app1 在线 00天00小时00分钟03秒."
         },
         {
             "id": "app1/a6",
             "info": "节点:a6(10.0.0.254) app1 在线 00天00小时00分钟02秒."
         },
         {
             "id": "app1/a7",
             "info": "节点:a7(10.0.0.254) app1 在线 00天00小时00分钟02秒."
         },
         {
             "id": "app1/a8",
             "info": "节点:a8(10.0.0.254) app1 在线 00天00小时00分钟01秒."
         },
         {
             "id": "app1/a9",
             "info": "节点:a9(10.0.0.254) app1 在线 00天00小时00分钟01秒."
         }
     ],
     "msg": "success"
 }`







# 生成密钥

```bash
#!/bin/bash
lang=2048
days=36500
serverName="$RANDOM.offer365.io"
echo ${serverName} > server_name.txt
subj="/C=CN/ST=BJ/L=BJ/O=offer365/OU=odin/CN=${serverName}"
# 为服务器和客户端分别生成私钥和证书
openssl genrsa -out server.key ${lang}
openssl req -new -x509 -days ${days} -subj "${subj}" -key server.key -out server.crt

openssl genrsa -out client.key ${lang}
openssl req -new -x509 -days ${days} -subj "${subj}" -key client.key -out client.crt

# 生成根证书
openssl genrsa -out ca.key ${lang}
openssl req -new -x509 -days ${days} -subj "${subj}" -key ca.key -out ca.crt

# 重新对服务器端证书进行签名
openssl req -new -subj "${subj}" -key server.key -out server.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days ${days} -in server.csr -out server.crt


## 重新对客户端证书进行签名
openssl req -new -subj "${subj}" -key client.key -out client.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days ${days} -in client.csr -out client.crt
```

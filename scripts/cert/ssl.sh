#!/bin/bash
# 涓烘湇鍔″櫒鍜屽鎴风鍒嗗埆鐢熸垚绉侀挜鍜岃瘉涔�
lang=2048
days=36500
openssl genrsa -out server.key ${lang}
openssl req -new -x509 -days ${days} -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" -key server.key -out server.crt

openssl genrsa -out client.key ${lang}
openssl req -new -x509 -days ${days} -subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io" -key client.key -out client.crt

# 鐢熸垚鏍硅瘉涔�
openssl genrsa -out ca.key ${lang}
openssl req -new -x509 -days ${days} -subj "/C=GB/L=China/O=gobook/CN=github.com" -key ca.key -out ca.crt

# 閲嶆柊瀵规湇鍔″櫒绔瘉涔﹁繘琛岀鍚�
openssl req -new -subj "/C=GB/L=China/O=server/CN=server.io" -key server.key -out server.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days ${days} -in server.csr -out server.crt

## 閲嶆柊瀵瑰鎴风璇佷功杩涜绛惧悕
openssl req -new \-subj "/C=GB/L=China/O=client/CN=client.io" -key client.key -out client.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days ${days} -in client.csr -out client.crt

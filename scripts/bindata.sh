# 生成静态文件 在本项目的根目录执行
go-bindata -o=asset/asset.go -pkg=asset html/... static/...

#go get -u github.com/mjibson/esc
#esc -pkg asset -o asset/static.go static  html

# 生成自签名的ca证书
openssl genrsa -out key.pem 2048
openssl req -new -x509 -key key.pem -out cert.pem -days 36500

# 安装 生成证书的工具
go get -u github.com/FiloSottile/mkcert
mkcert -install
mkcert odin #生成证书

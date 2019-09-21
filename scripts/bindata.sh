# 閻㈢喐鍨氶棃娆愶拷浣规瀮娴狅拷 閸︺劍婀版い鍦窗閻ㄥ嫭鐗撮惄顔肩秿閹笛嗩攽
go-bindata -o=asset/asset.go -pkg=asset html/... static/...

#go get -u github.com/mjibson/esc
#esc -pkg asset -o asset/static.go static  html

# 閻㈢喐鍨氶懛顏嗩劮閸氬秶娈慶a鐠囦椒鍔�
openssl genrsa -out key.pem 2048
openssl req -new -x509 -key key.pem -out cert.pem -days 36500

# 鐎瑰顥� 閻㈢喐鍨氱拠浣峰姛閻ㄥ嫬浼愰崗锟�
go get -u github.com/FiloSottile/mkcert
mkcert -install
mkcert odin #閻㈢喐鍨氱拠浣峰姛

# 闂佹眹鍨婚崰鎰板垂濮樿埖顥堟繛鍡樺姧閹烽攱鎷呯憴鍕拷顔济归悪鍛 闂侀潻璐熼崝宥咃耿閻楀牄浜滈柛锔诲幗缁愭鏌ｉ妸銉ヮ仾闁绘鎸抽幆鍕敊閼测晝协闂佸湱鐟抽崱鈺傛杸
go-bindata -o=asset/asset.go -pkg=asset html/... static/...

#go get -u github.com/mjibson/esc
#esc -pkg asset -o asset/static.go static  html

# 闂佹眹鍨婚崰鎰板垂濮樿埖鍤婃い蹇撴－閸旑噣鏌涘顒傂ゆ繛鍫熷従a闁荤姴娲ｅ鎺楀礉閿燂拷
openssl genrsa -out key.pem 2048
openssl req -new -x509 -key key.pem -out cert.pem -days 36500

# 闁诲海鎳撻ˇ鎶剿夐敓锟� 闂佹眹鍨婚崰鎰板垂濮樿鲸瀚氬ù锝呭槻婵盯鏌ｉ妸銉ヮ仼濞村吋鍔欏畷妤呮晸閿燂拷
go get -u github.com/FiloSottile/mkcert
mkcert -install
mkcert odin #闂佹眹鍨婚崰鎰板垂濮樿鲸瀚氬ù锝呭槻婵拷

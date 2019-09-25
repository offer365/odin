# 闂備焦鐪归崹濠氬窗閹版澘鍨傛慨妯垮煐椤ュ牊绻涢崱妯哄Ё闁圭兘鏀遍幏鍛喆閸曨叏鎷烽娴庡綊鎮崨顔碱伓 闂備線娼荤拹鐔煎礉瀹ュ拑鑰块柣妤�鐗勬禍婊堟煕閿旇骞楃紒鎰殜閺岋綁濡搁妷銉痪闂佺粯顨嗛幐鎶藉箚閸曨垼鏁婇柤娴嬫櫇鍗忛梻浣告贡閻熸娊宕遍埡鍌涙澑
go-bindata -o=asset/asset.go -pkg=asset html/... static/...

#go get -u github.com/mjibson/esc
#esc -pkg asset -o asset/static.go static  html

# 闂備焦鐪归崹濠氬窗閹版澘鍨傛慨妯垮煐閸ゅ﹥銇勮箛鎾达紞闁告棏鍣ｉ弻娑橆潩椤掑倐銈嗙箾閸喎寰揳闂佽崵濮村ú锝咁渻閹烘绀夐柨鐕傛嫹
openssl genrsa -out key.pem 2048
openssl req -new -x509 -key key.pem -out cert.pem -days 36500

# 闂佽娴烽幊鎾凰囬幎鍓垮鏁撻敓锟� 闂備焦鐪归崹濠氬窗閹版澘鍨傛慨妯块哺鐎氭艾霉閿濆懎妲诲┑顔界洴閺岋綁濡搁妷銉患婵炴潙鍚嬮崝娆忕暦濡ゅ懏鏅搁柨鐕傛嫹
go get -u github.com/FiloSottile/mkcert
mkcert -install
mkcert odin #闂備焦鐪归崹濠氬窗閹版澘鍨傛慨妯块哺鐎氭艾霉閿濆懎妲诲┑顕嗘嫹

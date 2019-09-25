#!/bin/bash
# appctl.sh
# 闂備礁鎼悧鐐哄箯閻戣姤鐓曟慨姗嗗墮椤掋垺绻涢幓鎺撳仴鐎殿噮鍓熼幃褔宕奸悤浣剐濋梻鍌氬�搁崑濠囧箯閿燂拷: 2019-08-04
pwd=123
ip=127.0.0.1
port=8888
address="http://${ip}:${port}/odin/api/v1"

function usage() {
  echo -e "
    getlic    闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉墮缁犳娊鏌涢銈呮瀾缂傚稄鎷峰┑鐑囩到濞茬娀宕滃┑鍥ь嚤闁跨噦鎷�
    putlic key    闂佽娴烽弫鎼佸储瑜斿畷锝夊幢濞戞顔婇梺鍓插亞閸犳劗绱為敓浠嬫⒑濞茶娅欓柟鍑ゆ嫹
    getcode    闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉ｅ妽閸庣喖鐓崶銊﹀碍闁诲繑鐟╅弻娑樷枎濞嗘劕顏�
    resetcode    闂傚倷鐒﹁ぐ鍐矓閸洘鍋柛鈩兩戦崕鐔肩叓閸ャ劍灏柣蹇旂懇閺屾稑鈻庡▎鎰伓
    qrcode    闂佺懓鍚嬪娆戞崲閹版澘鍨傛い蹇撶墕閻鏌у顒�锟藉鎯傚畝锟界槐鎾存媴閸撴彃娈堕梺缁樼壄閹凤拷
    dellic    婵犵數鍋涢ˇ顓㈠礉閹达箑缁╅柕蹇嬪�曠粻鎶芥煕椤愩倕鏋戠紓宥忔嫹
    qrlic    婵犵數鍋涢ˇ顓㈠礉閹达箑缁╅柕蹇ョ磿椤╄尙锟界櫢鎷烽柨鐔绘鑿愭い鏇嫹闁诡垰绉归弫鎾绘晸閿燂拷
    nodes    闂備胶鍘ч幖顐﹀磹婵犳艾纾婚柨婵嗘川绾惧ジ鏌涢弴銊ュ箻缂佸鎷�
    conf    闂傚倷鐒﹀妯肩矓閸洘鍋柛鈩冾焽閳瑰秹鏌嶉埡浣告殨缂佹唻鎷�
    online  闂佽楠哥粻宥夊垂濞差亜鏄ユ繛鎴炴皑閸楁碍銇勯弽銊ь暡闂婎剦鍓涚槐鎺楀箻閸涱喖顏�"
}

# 闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉墮缁犳娊鏌涢銈呮瀾缂傚稄鎷峰┑鐑囩到濞茬娀宕滃┑鍥ь嚤闁跨噦鎷�
function getlic() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# 闂佽娴烽弫鎼佸储瑜斿畷锝夊幢濞戞顔婇梺鍓插亞閸犳劗绱為敓锟�
function putlic() {
  curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# 闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉ｅ妽閸庣喖鐓崶銊﹀碍闁诲繑鐟╅弻娑樷枎濞嗘劕顏�
function getcode() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/code"
}
# 闂傚倷鐒﹁ぐ鍐矓閸洘鍋柛鈩兩戦崕鐔肩叓閸ャ劍灏柣蹇旂懇閺屾稑鈻庡▎鎰伓
function resetcode() {
  curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# 闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉ｅ妽閸庣喖鐓崶銊﹀碍闁诲繑鐟╅弻娑樷枎閹烘捁鍚傜紓浣靛妷閺呮繄妲愰幒鏃撴嫹闂堟稒鎲搁柛顭掓嫹
function qrcode() {
  curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# 婵犵數鍋涢ˇ顓㈠礉閹达箑缁╅柕蹇嬪�曠粻鎶芥煕椤愩倕鏋戠紓宥忔嫹
function deletelic() {
  curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# 闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉戯拷閺嬫牠鏌曟繛鐐珕闁轰焦锚闇夋繝濠傚暟閻忚鲸绻涢敐搴℃珝闁诡垰绉归弫鎾绘晸閿燂拷
function qrlicense() {
  curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# 闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉墯閸ゅ秹鏌涚仦鍓х煀濞寸媭鍨拌彁闁搞儻绲芥晶鎻捗圭涵閿嬪
function nodes() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# 闂備礁鍚嬮崕鎶藉床閼艰翰浜归柛銉墯閻擄拷閻庡箍鍎卞ú銊╁几閸屾搴ㄥ炊閿濆懎鈷夋繛瀵搞�嬮幏锟�
function conf() {
  curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# 闂佽楠哥粻宥夊垂濞差亜鏄ユ繛鎴炴皑閸楁碍銇勯弽銊ь暡闂婎剦鍓涚槐鎺楀箚瑜嬫禒鎺懬庨崶褝韬鐐扮窔閺佹捇鏁撻敓锟�
function online() {
  curl -k -s -X GET --user admin:${pwd} "${address}/client/online/$1"
}

case "$1" in
"getlic")
  getlic
  ;;
"putlic")
  putlic $2
  ;;
"getcode")
  getcode
  ;;
"resetcode")
  resetcode
  ;;
"qrcode")
  qrcode
  ;;
"dellic")
  deletelic
  ;;
"qrlic")
  qrlicense
  ;;
"nodes")
  nodes
  ;;
"conf")
  conf $2
  ;;
"online")
  online $2
  ;;
*)
  usage
  ;;
esac

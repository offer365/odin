#!/bin/bash
# appctl.sh
# 闂佸搫鐗為幏鐑芥煕濮橆剙顒㈡繛鎻掓健瀵剟鎮ч崼鐕佹Н闂傚倸鍋婇幏锟�: 2019-08-04
pwd=123
ip=127.0.0.1
port=8888
address="http://${ip}:${port}/odin/api/v1"

function usage() {
  echo -e "
    getlic    闂佸吋鍎抽崲鑼躲亹閸ヮ剙绠抽柛顐ゅ枑缂嶏拷婵烇絽娲犻崜婵囧閿燂拷
    putlic key    闁诲海鏁搁崢褔宕ｉ崱娑樼闁割偆鍠愮紞锟介梺娲诲櫙閹凤拷
    getcode    闂佸吋鍎抽崲鑼躲亹閸ャ劍鍎熼煫鍥ㄦ尭閻忔瑩鏌涘▎娆愬
    resetcode    闂備焦褰冪粔鍫曟偪閸℃ɑ鍎熼煫鍥ㄦ尭閻忔瑩鏌涘▎娆愬
    qrcode    闁瑰吋娼欑换鎰板垂椤忓牆鐭楅柧姘�婚惂宀�绱撴担鍓插殶闁绘牭鎷�
    dellic    濠电偛顦崝鎴﹀绩閵忋倕绠抽柛顐ゅ枑缂嶏拷
    qrlic    濠电偛顦崝鎴﹀绩閵忥紕顩茬�癸拷閿熻棄菐椤曪拷閹秹鏁撻敓锟�
    nodes    闂佺厧鎼崐濠氬磻閿濆洨纾奸柛鏇ㄥ幘缁夛拷
    conf    闂備焦婢樼粔鍫曟偪閸℃鈹嶉柍鈺佸暕缁憋拷
    online  闁诲骸绠嶉崹娲春濞戞氨鍗氭い鏍ㄧ闊剛绱掗幘鍛"
}

# 闂佸吋鍎抽崲鑼躲亹閸ヮ剙绠抽柛顐ゅ枑缂嶏拷婵烇絽娲犻崜婵囧閿燂拷
function getlic() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# 闁诲海鏁搁崢褔宕ｉ崱娑樼闁割偆鍠愮紞锟�
function putlic() {
  curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# 闂佸吋鍎抽崲鑼躲亹閸ャ劍鍎熼煫鍥ㄦ尭閻忔瑩鏌涘▎娆愬
function getcode() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/code"
}
# 闂備焦褰冪粔鍫曟偪閸℃ɑ鍎熼煫鍥ㄦ尭閻忔瑩鏌涘▎娆愬
function resetcode() {
  curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# 闂佸吋鍎抽崲鑼躲亹閸ャ劍鍎熼煫鍥ㄦ尭閻忔瑩鏌涘▎鎺撹吂缂併劉鏅濈槐鎺旓拷闈涙憸閸拷
function qrcode() {
  curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# 濠电偛顦崝鎴﹀绩閵忋倕绠抽柛顐ゅ枑缂嶏拷
function deletelic() {
  curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# 闂佸吋鍎抽崲鑼躲亹閸パ�鏋栭柕濞炬櫆閺佹ê霉濠婂啰鐏辨繛锝庡櫍閹秹鏁撻敓锟�
function qrlicense() {
  curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# 闂佸吋鍎抽崲鑼躲亹閸ヮ剚鍤嶉柛灞剧矊娴狀垰菐閸ワ絽澧插ù纭锋嫹
function nodes() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# 闂佸吋鍎抽崲鑼躲亹閸ヮ剚鐓�鐎广儱娲ㄩ弸鍌毲庨崶锝呭⒉濞寸》鎷�
function conf() {
  curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# 闁诲骸绠嶉崹娲春濞戞氨鍗氭い鏍ㄧ闊剛绱掗幆褋浠掑ǎ鍥э躬楠炰線鏁撻敓锟�
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

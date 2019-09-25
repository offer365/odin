#!/bin/bash
# appctl.sh
# 闁哄牞鎷烽柛姘濞插潡寮悧鍫燁槯闂傚偊鎷�: 2019-08-04
pwd=123
ip=127.0.0.1
port=8888
address="http://${ip}:${port}/odin/api/v1"

function usage() {
  echo -e "
    getlic    闁兼儳鍢茶ぐ鍥箳閸喐缍�濞ｅ洠鍓濇导锟�
    putlic key    閻庣數鍘ч崣鍡涘箳閸喐缍�闁活噯鎷�
    getcode    闁兼儳鍢茶ぐ鍥ㄦ償韫囨挸鐏欓柛娆欐嫹
    resetcode    闂佹彃绉堕悿鍡樻償韫囨挸鐏欓柛娆欐嫹
    qrcode    閹兼潙绻愰崹顏堝矗閾氬倻鐧岀紓浣割嚟閻栵拷
    dellic    婵炲鍔戦弨銏ゅ箳閸喐缍�
    qrlic    婵炲鍔戦弨銏＄瀹�锟藉ǎ顕�鎯嶉敓锟�
    nodes    闁煎搫鍊婚崑锝囩磼閸曨厾绉�
    conf    闂佹澘绉堕悿鍡樼┍閳╁啩绱�
    online  閻庡箍鍨洪崺娑氱博椤栨碍韬紒鎾呮嫹"
}

# 闁兼儳鍢茶ぐ鍥箳閸喐缍�濞ｅ洠鍓濇导锟�
function getlic() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# 閻庣數鍘ч崣鍡涘箳閸喐缍�
function putlic() {
  curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# 闁兼儳鍢茶ぐ鍥ㄦ償韫囨挸鐏欓柛娆欐嫹
function getcode() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/code"
}
# 闂佹彃绉堕悿鍡樻償韫囨挸鐏欓柛娆欐嫹
function resetcode() {
  curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# 闁兼儳鍢茶ぐ鍥ㄦ償韫囨挸鐏欓柛娆掓腹缁ㄢ晝绱掔�靛摜鍨�
function qrcode() {
  curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# 婵炲鍔戦弨銏ゅ箳閸喐缍�
function deletelic() {
  curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# 闁兼儳鍢茶ぐ鍥р枖閵娾晜鏁樺ù婊冪灱濞ｎ噣鎯嶉敓锟�
function qrlicense() {
  curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# 闁兼儳鍢茶ぐ鍥嚍閸屾粌浠ǎ鍥ｅ墲娴硷拷
function nodes() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# 闁兼儳鍢茶ぐ鍥煀瀹ュ洨鏋傚ǎ鍥ｅ墲娴硷拷
function conf() {
  curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# 閻庡箍鍨洪崺娑氱博椤栨碍韬紒鎯с仒娣囧﹪骞侀敓锟�
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

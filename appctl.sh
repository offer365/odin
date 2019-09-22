#!/bin/bash
# appctl.sh
# 閺堬拷閸氬孩娲块弬鐗堟闂傦拷: 2019-08-04
pwd=123
ip=127.0.0.1
port=8888
address="http://${ip}:${port}/odin/api/v1"

function usage() {
  echo -e "
    getlic    閼惧嘲褰囬幒鍫熸綀娣団剝浼�
    putlic key    鐎电厧鍙嗛幒鍫熸綀閻拷
    getcode    閼惧嘲褰囨惔蹇撳灙閸欙拷
    resetcode    闁插秶鐤嗘惔蹇撳灙閸欙拷
    qrcode    鎼村繐鍨崣铚傜癌缂佸鐖�
    dellic    濞夈劑鏀㈤幒鍫熸綀
    qrlic    濞夈劑鏀㈡禍宀�娣惍锟�
    nodes    閼哄倻鍋ｇ紒鍕秹
    conf    闁板秶鐤嗘穱鈩冧紖
    online  鐎广垺鍩涚粩顖氭躬缁撅拷"
}

# 閼惧嘲褰囬幒鍫熸綀娣団剝浼�
function getlic() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# 鐎电厧鍙嗛幒鍫熸綀
function putlic() {
  curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# 閼惧嘲褰囨惔蹇撳灙閸欙拷
function getcode() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/code"
}
# 闁插秶鐤嗘惔蹇撳灙閸欙拷
function resetcode() {
  curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# 閼惧嘲褰囨惔蹇撳灙閸欒渹绨╃紒瀵哥垳
function qrcode() {
  curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# 濞夈劑鏀㈤幒鍫熸綀
function deletelic() {
  curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# 閼惧嘲褰囧▔銊╂敘娴滃瞼娣惍锟�
function qrlicense() {
  curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# 閼惧嘲褰囬懞鍌滃仯娣団剝浼�
function nodes() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# 閼惧嘲褰囬柊宥囩枂娣団剝浼�
function conf() {
  curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# 鐎广垺鍩涚粩顖氭躬缁惧じ淇婇幁锟�
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

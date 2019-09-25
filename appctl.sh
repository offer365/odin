#!/bin/bash
# appctl.sh
# 最后更新时间: 2019-08-04
pwd=123
ip=127.0.0.1
port=8888
address="http://${ip}:${port}/odin/api/v1"


function usage() {
    echo -e "
    getlic    获取授权信息
    putlic key    导入授权码
    getcode    获取序列号
    resetcode    重置序列号
    qrcode    序列号二维码
    dellic    注销授权
    qrlic    注销二维码
    nodes    节点组网
    conf    配置信息
    online  客户端在线"
}

# 获取授权信息
function getlic() {
    curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# 导入授权
function putlic() {
    curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# 获取序列号
function getcode() {
    curl -k -s -X GET --user admin:${pwd}  "${address}/server/code"
}
# 重置序列号
function resetcode() {
    curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# 获取序列号二维码
function qrcode() {
    curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# 注销授权
function deletelic() {
    curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# 获取注销二维码
function qrlicense() {
    curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# 获取节点信息
function nodes() {
    curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# 获取配置信息
function conf() {
    curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# 客户端在线信息
function online() {
    curl -k -s -X GET --user admin:${pwd} "${address}/client/online/$1"
}

case "$1" in
    "getlic")
        getlic ;;
    "putlic")
        putlic $2 ;;
    "getcode")
        getcode ;;
    "resetcode")
        resetcode ;;
    "qrcode")
        qrcode ;;
    "dellic")
        deletelic ;;
    "qrlic")
        qrlicense ;;
    "nodes")
        nodes ;;
    "conf")
        conf $2 ;;
    "online")
        online $2;;
    *)
        usage ;;
esac

#!/bin/bash
# appctl.sh
# 最后更新时间: 2019-08-04
pwd=123
ip=10.0.0.200
port=9527
address="http://${ip}:${port}/odin/api/v1"


function usage() {
    echo -e "
    getlic
    putlic key
    getcode
    resetcode
    qrcode
    dellic
    qrlic
    nodes
    conf
    online  "
}

# get authorization info
function getlic() {
    curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# immport authorization
function putlic() {
    curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# get serial number
function getcode() {
    curl -k -s -X GET --user admin:${pwd}  "${address}/server/code"
}
# Reset serial number
function resetcode() {
    curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# Get the serial number QR code
function qrcode() {
    curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# Logout authorization
function deletelic() {
    curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# Get the logout QR code
function qrlicense() {
    curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# get node configuration
function nodes() {
    curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# get config info
function conf() {
    curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# client onlien info
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
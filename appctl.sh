#!/bin/bash
# appctl.sh
# 闂傚倷绀侀幖顐︽偋閻愬搫绠柣鎴ｅГ閻撴洘鎱ㄥ鍡楀妞ゆ帇鍨虹换娑㈠箵閹烘挸浠撮悗娈垮櫘閸撶喖骞冭瀹曞ジ鎮ゆ担鍓愭繈姊婚崒姘拷鎼佸磻婵犲洤绠柨鐕傛嫹: 2019-08-04
pwd=123
ip=127.0.0.1
port=8888
address="http://${ip}:${port}/odin/api/v1"

function usage() {
  echo -e "
    getlic    闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷顔煎缂佺姵濞婇弻娑㈩敃閵堝懏鐎剧紓鍌氱▌閹峰嘲鈹戦悜鍥╁埌婵炶尙濞�瀹曟粌鈹戦崶褜鍤ら梺璺ㄥ櫐閹凤拷
    putlic key    闂備浇顕уù鐑藉极閹间礁鍌ㄧ憸鏂跨暦閿濆骞㈡繛鎴烆焽椤斿﹪姊洪崜鎻掍簽闁哥姵鍔楃槐鐐烘晸娴犲鈷戞繛鑼额嚙濞呮瑩鏌熼崙銈嗗
    getcode    闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷锝呭闁稿海鍠栭悡顐﹀炊閵婏箑纰嶉梺璇茬箲閻熲晠寮诲☉妯锋瀻婵炲棙鍔曢锟�
    resetcode    闂傚倸鍊烽悞锕併亹閸愵亞鐭撻柛顐ｆ礃閸嬵亪鏌涢埄鍏╂垿宕曢悢鑲╁彄闁搞儯鍔嶇亸顓㈡煟韫囨梻鎳囬柡灞剧☉閳诲骸鈻庨幇顒変紦
    qrcode    闂備胶鎳撻崥瀣焽濞嗘垶宕查柟鐗堟緲閸ㄥ倹銇勮箛鎾跺闁活厽顨婇弻褍顫濋锟介敓钘夘煼閹倸鐣濋敓鐣屾閹惧瓨濯撮柛鎾村絻濞堝爼姊虹紒妯煎闁瑰嚖鎷�
    dellic    濠电姷鏁搁崑娑⑺囬銏犵闁硅揪绠戠紒鈺呮煏韫囧锟芥洜绮婚幎鑺ョ厱妞ゆ劑鍊曢弸鎴犵磽瀹ュ繑瀚�
    qrlic    濠电姷鏁搁崑娑⑺囬銏犵闁硅揪绠戠紒鈺呮煏韫囥儳纾挎い鈺勫皺閿熺晫娅㈤幏鐑芥煥閻旂粯顥夐懣鎰亜閺囶亝瀚归梺璇″灠缁夊綊寮幘缁樻櫢闁跨噦鎷�
    nodes    闂傚倷鑳堕崢褔骞栭锕�纾瑰┑鐘宠壘绾惧鏌ㄥ┑鍡樺窛缁炬儳銈搁弻娑㈠即閵娿儱绠荤紓浣割檧閹凤拷
    conf    闂傚倸鍊烽悞锕�顭垮Ο鑲╃煋闁割偅娲橀崑顏堟煕閳╁喚鐒介柍鐟扮Ч閺屽秹鍩℃担鍛婃缂備焦鍞婚幏锟�
    online  闂備浇顕ф鍝ョ不瀹ュ鍨傛繛宸簻閺勩儲绻涢幋鐐寸殤闁告纰嶉妵鍕冀閵娧屾殹闂傚鍓﹂崜娑氭閹烘绠婚柛娑卞枛椤忥拷"
}

# 闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷顔煎缂佺姵濞婇弻娑㈩敃閵堝懏鐎剧紓鍌氱▌閹峰嘲鈹戦悜鍥╁埌婵炶尙濞�瀹曟粌鈹戦崶褜鍤ら梺璺ㄥ櫐閹凤拷
function getlic() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/license"
}
# 闂備浇顕уù鐑藉极閹间礁鍌ㄧ憸鏂跨暦閿濆骞㈡繛鎴烆焽椤斿﹪姊洪崜鎻掍簽闁哥姵鍔楃槐鐐烘晸閿燂拷
function putlic() {
  curl -k -s -X POST --user admin:${pwd} -F key="$1" "${address}/server/license"
}
# 闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷锝呭闁稿海鍠栭悡顐﹀炊閵婏箑纰嶉梺璇茬箲閻熲晠寮诲☉妯锋瀻婵炲棙鍔曢锟�
function getcode() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/code"
}
# 闂傚倸鍊烽悞锕併亹閸愵亞鐭撻柛顐ｆ礃閸嬵亪鏌涢埄鍏╂垿宕曢悢鑲╁彄闁搞儯鍔嶇亸顓㈡煟韫囨梻鎳囬柡灞剧☉閳诲骸鈻庨幇顒変紦
function resetcode() {
  curl -k -s -X POST --user admin:${pwd} "${address}/server/code"
}
# 闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷锝呭闁稿海鍠栭悡顐﹀炊閵婏箑纰嶉梺璇茬箲閻熲晠寮诲☉妯锋瀻闁圭儤鎹侀崥鍌滅磽娴ｉ潧濡烽柡鍛箘濡叉劙骞掗弮鎾村闂傚牊绋掗幉鎼佹煕椤帗瀚�
function qrcode() {
  curl -k -s -X GET --user admin:${pwd} -o qr-code.jpg "${address}/server/qr-code"
}
# 濠电姷鏁搁崑娑⑺囬銏犵闁硅揪绠戠紒鈺呮煏韫囧锟芥洜绮婚幎鑺ョ厱妞ゆ劑鍊曢弸鎴犵磽瀹ュ繑瀚�
function deletelic() {
  curl -k -s -X DELETE --user admin:${pwd} "${address}/server/license"
}
# 闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷鎴嫹闁哄鐗犻弻鏇熺箾閻愵剚鐝曢梺杞扮劍閿氶棁澶嬬節婵犲倸鏆熼柣蹇氶哺缁绘盯鏁愭惔鈩冪彎闂佽鍨扮粔褰掑极閹剧粯鏅搁柨鐕傛嫹
function qrlicense() {
  curl -k -s -X GET --user admin:${pwd} -o qr-license.jpg "${address}/server/qr-license"
}

# 闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷顔煎闁搞倕绉归弻娑氫沪閸撗呯厐婵炲濯崹鎷屽絹闂佹悶鍎荤徊鑺ユ櫠閹绘崡鍦兜闁垮顏�
function nodes() {
  curl -k -s -X GET --user admin:${pwd} "${address}/server/nodes"
}
# 闂傚倷绀侀崥瀣磿閹惰棄搴婇柤鑹扮堪娴滃綊鏌涢妷顔煎闁绘搫鎷烽柣搴＄畭閸庡崬煤閵娾晛鍑犻柛灞绢嚙鎼淬劌鐐婇柨婵嗘噹閳峰绻涚�垫悶锟藉骞忛敓锟�
function conf() {
  curl -k -s -X GET --user admin:${pwd} "${address}/client/conf/$1"
}
# 闂備浇顕ф鍝ョ不瀹ュ鍨傛繛宸簻閺勩儲绻涢幋鐐寸殤闁告纰嶉妵鍕冀閵娧屾殹闂傚鍓﹂崜娑氭閹烘绠氱憸瀣閹烘嚞搴ㄥ炊瑜濋煬顒侇殽閻愭壆绐旈柡浣规崌閺佹捇鏁撻敓锟�
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

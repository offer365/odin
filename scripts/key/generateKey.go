package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/offer365/example/endecrypt/endeecc"
	"github.com/offer365/example/endecrypt/endersa"
	"github.com/offer365/example/tools"
)

// 为odin 生成各种秘钥

const (
	temp = `package main

const  (
	_eccpri1 = #{_eccpri1}#
	_eccpub1 = #{_eccpub1}#
	_eccpri2 = #{_eccpri2}#
	_eccpub2 = #{_eccpub2}#
	_eccpri3 = #{_eccpri3}#
	_eccpub3 = #{_eccpub3}#

	_rsa1024pri1 = #{_rsa1024pri1}#
	_rsa1024pub1 = #{_rsa1024pub1}#
	_rsa1024pri2 = #{_rsa1024pri2}#
	_rsa1024pub2 = #{_rsa1024pub2}#
	_rsa1024pri3 = #{_rsa1024pri3}#
	_rsa1024pub3 = #{_rsa1024pub3}#
	_rsa2048pri1 = #{_rsa2048pri1}#
	_rsa2048pub1 = #{_rsa2048pub1}#
	_rsa2048pri2 = #{_rsa2048pri2}#
	_rsa2048pub2 = #{_rsa2048pub2}#
	_rsa2048pri3 = #{_rsa2048pri3}#
	_rsa2048pub3 = #{_rsa2048pub3}#
	_rsa4096pri1 = #{_rsa4096pri1}#
	_rsa4096pub1 = #{_rsa4096pub1}#
	_rsa4096pri2 = #{_rsa4096pri2}#
	_rsa4096pub2 = #{_rsa4096pub2}#
	_rsa4096pri3 = #{_rsa4096pri3}#
	_rsa4096pub3 = #{_rsa4096pub3}#

	_aes256key1 = #{_aes256key1}#
	_aes256key2 = #{_aes256key2}#
	_aes256key3 = #{_aes256key3}#
	_aes256key4 = #{_aes256key4}#
)
`
)

func main() {
	endersa.GetRsaKey(1024, "_rsa1024pri1.pem", "_rsa1024pub1.pem")
	endersa.GetRsaKey(1024, "_rsa1024pri2.pem", "_rsa1024pub2.pem")
	endersa.GetRsaKey(1024, "_rsa1024pri3.pem", "_rsa1024pub3.pem")

	endersa.GetRsaKey(2048, "_rsa2048pri1.pem", "_rsa2048pub1.pem")
	endersa.GetRsaKey(2048, "_rsa2048pri2.pem", "_rsa2048pub2.pem")
	endersa.GetRsaKey(2048, "_rsa2048pri3.pem", "_rsa2048pub3.pem")

	endersa.GetRsaKey(4096, "_rsa4096pri1.pem", "_rsa4096pub1.pem")
	endersa.GetRsaKey(4096, "_rsa4096pri2.pem", "_rsa4096pub2.pem")
	endersa.GetRsaKey(4096, "_rsa4096pri3.pem", "_rsa4096pub3.pem")

	endeecc.GetEccKey("_eccpri1.pem", "_eccpub1.pem")
	endeecc.GetEccKey("_eccpri2.pem", "_eccpub2.pem")
	endeecc.GetEccKey("_eccpri3.pem", "_eccpub3.pem")

	for i := 1; i < 5; i++ {
		key := tools.RandString(32)
		f, err := os.Create("_aes256key" + strconv.Itoa(i) + ".pem")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = f.Write([]byte(key))
		if err != nil {
			fmt.Println(err)
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	time.Sleep(time.Second * 2)
	str := temp
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		name := file.Name()
		data, err := ioutil.ReadFile(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		key := strings.Split(name, ".")[0]
		str = strings.Replace(str, "#{"+key+"}#", "`"+string(data)+"`", 1)
	}
	f, err := os.Create("key.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

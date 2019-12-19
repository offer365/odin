package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/offer365/example/tools"
)

// 为odin 生成各种秘钥

const (
	temp = `package main

const  (
	server_crt=#{server_crt}#
    server_csr=#{server_csr}#
	server_key=#{server_key}#
	client_crt=#{client_crt}#
    client_csr=#{client_csr}#
	client_key=#{client_key}#
	ca_crt=#{ca_crt}#
	ca_key=#{ca_key}#
	ca_srl=#{ca_srl}#
    server_name=#{server_name}#

	clusterToken = "#{sha256sum}#"
	embedAuthPwd = "#{sha256sum}#"

	storeLicenseKey            = "/#{sha256sum}#/#{sha256sum}#"
	storeClearLicenseKey       = "/#{sha256sum}#/#{sha256sum}#"
	storeClientConfigKeyPrefix = "/#{sha256sum}#/#{sha256sum}#/"
	storeClientKeyPrefix       = "/#{sha256sum}#/#{sha256sum}#/"
	storeTokenKey              = "/#{sha256sum}#/#{sha256sum}#"
	storeSerialNumKey          = "/#{sha256sum}#/#{sha256sum}#"
	storeHashSalt              = "#{sha256sum}#"

	grpcUser      = "#{sha256sum}#"
	grpcPwd       = "#{sha256sum}#"
)
`
)

func main() {
	str:=temp
	for strings.Contains(str,"#{sha256sum}#"){
		s:=tools.Sha256Hex([]byte(tools.RandString(16)),nil)
		str=strings.Replace(str,"#{sha256sum}#",s,1)
	}


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
		keys := strings.Split(name, ".")
		str = strings.Replace(str, "#{"+strings.Join(keys,"_")+"}#", "`"+string(data)+"`", 1)
	}
	f, err := os.Create("ssl.go")
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




package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

// 运行时间
func RunTime(now, start int64) string {
	if now < start {
		return ""
	}
	online := now - start
	d := online / 86400
	h := (online - d*86400) / 3600
	m := (online - d*86400 - h*3600) / 60
	s := online - d*86400 - h*3600 - m*60
	return fmt.Sprintf("%02d天%02d小时%02d分钟%02d秒.", d, h, m, s)
}

// 绝对值
func Abs(a int64) int64 {
	return (a ^ a>>31) - a>>31
}

func Md5sum(byt, salt []byte) string {
	h := md5.New()
	if salt != nil {
		byt = append(byt, salt...)
	}
	h.Write(byt)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Sha256sum(byt, salt []byte) string {
	h := sha256.New()
	if salt != nil {
		byt = append(byt, salt...)
	}
	h.Write(byt)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Scrypt(src, salt []byte) string {
	byt,err:=scrypt.Key(src, salt, 1<<15, 8, 1, 32)
	if err!=nil{
		return ""
	}
	return base64.StdEncoding.EncodeToString(byt)
}

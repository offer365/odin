package dao

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	mmm()
}

func mmm() {
	h := md5.New()
	h.Write([]byte("aa"))
	a := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(a)
}

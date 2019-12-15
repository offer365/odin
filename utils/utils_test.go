package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestAbs(t *testing.T) {
	fmt.Println(Abs(-333))
}

func TestRunTime(t *testing.T) {
	fmt.Println(RunTime(time.Now().Unix()+784521, time.Now().Unix()))
}

func TestMd5sum(t *testing.T) {
	fmt.Println(Md5sum([]byte("123"), nil))
	fmt.Println(Md5sum([]byte("123"), []byte("456")))
}

func TestSha256sum(t *testing.T) {
	fmt.Println(Sha256sum([]byte("123"), nil))
	fmt.Println(Sha256sum([]byte("123"), []byte("456")))
}

package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHw_Md5(t *testing.T) {
	h := NewNode([]string{"192.168.10.101"})
	byt, _ := json.Marshal(h.Hardware)
	fmt.Println(string(byt))
	byt, _ = json.Marshal(h)
	fmt.Println(string(byt))
}

func TestGenSerialNum(t *testing.T) {

}

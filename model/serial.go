package model

import (
	"encoding/json"
	"github.com/offer365/endecrypt"
	"github.com/offer365/endecrypt/endeaesrsa"
	"github.com/offer365/odin/node"
	"github.com/satori/go.uuid"
	"time"
)

// 序列号
type SerialNum struct {
	Sid   string                `json:"sid"`   // 序列号唯一uuid，用来标识序列号，并与 授权码相互校验，一一对应。
	Nodes map[string]*node.Node `json:"nodes"` // 节点的具体硬件信息。这里不使用map的原因是map是无序的。无法保证每次生成的hws是一致的。
	Time  int64                 `json:"time"`  // 生成 序列号的时间。
}

//生成序列号
func (sn *SerialNum) GenSerialNum(nodes map[string]*node.Node) (code string, err error) {
	var byt []byte
	sn.Nodes = nodes
	sn.Sid = uuid.Must(uuid.NewV4()).String()
	// 生成序列号的时间
	sn.Time = time.Now().Unix()
	// 序列化 实例
	if byt, err = json.Marshal(sn); err != nil {
		return
	}
	// 公钥加密 生成序列号
	return endeaesrsa.PubEncrypt(byt, endecrypt.PubkeyServer2048, endecrypt.AesKeyServer2)
}

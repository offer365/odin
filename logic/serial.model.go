package logic

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/offer365/example/endecrypt"
	pb "github.com/offer365/odin/proto"
	uuid "github.com/satori/go.uuid"
)

var (
	Serial *SerialNum
)

func init() {
	Serial = new(SerialNum)
}

// 序列号
type SerialNum struct {
	Sid   string              `json:"sid"`   // 序列号唯一uuid，用来标识序列号，并与 授权码相互校验，一一对应。
	Nodes map[string]*pb.Node `json:"nodes"` // 节点的具体硬件信息。
	Date  int64               `json:"date"`  // 生成 序列号的时间。
}

// 生成序列号
func (sn *SerialNum) Generate(nodes map[string]*pb.Node) (code string, err error) {
	var byt []byte
	sn.Nodes = nodes
	sn.Sid = uuid.Must(uuid.NewV4()).String()
	// 生成序列号的时间
	sn.Date = time.Now().Unix()
	// 序列化 实例
	if byt, err = json.Marshal(sn); err != nil {
		return
	}
	// 公钥加密 生成序列号
	byt, err = endecrypt.Encrypt(endecrypt.Pub1AesRsa2048, byt)
	if err != nil {
		return
	}
	code = base64.StdEncoding.EncodeToString(byt)
	return
}

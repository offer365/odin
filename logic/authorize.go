package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/offer365/endecrypt"
	"github.com/offer365/odin/proto"
	"github.com/offer365/odin/utils"
	uuid "github.com/satori/go.uuid"
)

var Author = &Authorize{}

type Authorize struct {
	sync.Mutex
}

const (
	statusOk        int32 = 200
	statusDecodeErr int32 = iota + 410
	statusDecryptErr
	statusUnmarshalErr
	statusValidatorErr
	statusCheckErr
	statusExpiresErr
	statusCountErr
	statusInsufficientErr
	statusExistsErr
	statusEncryptErr int32 = iota + 420
	statusPutClientErr
	statusMarshalErr
	statusPutTokenErr
	statusGetClientErr
	statusUmd5Err
	statusKeepErr
	statusOffErr
)

type Validator struct {
	App   string `json:"app"`
	ID    string `json:"id"`
	Date  int64  `json:"date"`
	Token string `json:"token"`
}

func (a *Authorize) Auth(ctx context.Context, req *proto.Request) (resp *proto.Response, err error) {
	var (
		byt             []byte
		app, id, token  string
		num             int64
		exist, register bool
	)
	// 解密
	if byt, err = base64.StdEncoding.DecodeString(req.Verify); err != nil {
		resp = &proto.Response{Code: statusDecodeErr, Msg: "decode string verify error"}
		return
	}

	if byt, err = endecrypt.Decrypt(endecrypt.Aes2key32, byt); err != nil {
		resp = &proto.Response{Code: statusDecryptErr, Msg: "decrypt verify error"}
		return
	}
	valid := new(Validator)
	if err = json.Unmarshal(byt, valid); err != nil {
		resp = &proto.Response{Code: statusUnmarshalErr, Msg: "unmarshal verify error"}
		return
	}
	// 检查请求是否合法
	if valid.App != valid.App || valid.ID != req.Id || valid.Date != req.Date || utils.Abs(time.Now().Unix()-valid.Date) > 600 {
		resp = &proto.Response{Code: statusValidatorErr, Msg: "verification failed"}
		err = errors.New("verification failed")
		return
	}
	app = valid.App
	id = valid.ID
	token = valid.Token
	a.Lock()
	defer a.Unlock()
	// token 是否存在 或者 不存在是否可以注册
	exist, register = GetTokenAndChk(app, id, token)
	// 1,token 不存在 不可注册 -----退出
	// 2,token 存在  -----下一步
	// 3,token 不存在 可注册 --- 下一步
	if !exist && !register {
		resp = &proto.Response{Code: statusCheckErr, Msg: "auth failed or token error"}
		return
	}
	// 检查应用是否授权到期
	if !LoadLic().CheckTime(app) {
		resp = &proto.Response{Code: statusExpiresErr, Msg: "app does not exist or authorization expires"}
		return
	}
	// 检查实例是否超出授权个数
	num, err = CountClient(app)
	if err != nil {
		resp = &proto.Response{Code: statusCountErr, Msg: "get the number of App instances"}
		return
	}
	if !LoadLic().ChkInstance(app, num) {
		resp = &proto.Response{Code: statusInsufficientErr, Msg: "app has insufficient remaining instances"}
		return
	}
	// 检查实例是否已经存在
	cli, exist := GetClient(app + "/" + id)
	if cli != nil || exist {
		resp = &proto.Response{Code: statusExistsErr, Msg: "the id already exists"}
		return
	}
	nc := &Cli{
		ID:    id,
		App:   app,
		Uuid:  uuid.Must(uuid.NewV4()).String(),
		Start: time.Now().Unix(),
		Lease: 0, // 租约id 在PutClient里面赋值
	}
	// 生成uuid密文
	cipher, err := endecrypt.Encrypt(endecrypt.Pri2Rsa1024, []byte(nc.Uuid))
	if err != nil {
		resp = &proto.Response{Code: statusEncryptErr, Msg: "encrypt uuid error"}
		return
	}
	// 生成10秒租约
	lease, err := PutClient(app+"/"+id, nc)
	if err != nil {
		resp = &proto.Response{Code: statusPutClientErr, Msg: "save instance error"}
		return
	}
	// 生成授权信息
	attrs := LoadLic().Apps[app].Attrs
	data := make(map[string]interface{})
	data["attrs"] = attrs
	data["time"] = time.Now().UnixNano() // 保证每次生成的attr 不一样
	byt, err = json.Marshal(data)
	if err != nil {
		resp = &proto.Response{Code: statusMarshalErr, Msg: "marshal authinfo error"}
		return
	}
	auth, err := endecrypt.Encrypt(endecrypt.Pri2Rsa2048, byt)
	if err != nil {
		resp = &proto.Response{Code: statusEncryptErr, Msg: "encrypt authinfo error"}
		return
	}

	if register {
		if err = PutToken(app, id, token); err != nil {
			resp = &proto.Response{Code: statusPutTokenErr, Msg: "put token error"}
			return
		}
	}

	// 生成的auth 与 cipher 使用不同的加密算法。
	resp = &proto.Response{
		Code: statusOk,
		Data: &proto.Data{
			Auth:   auth,
			Cipher: cipher,
			Lease:  lease,
		},
		Msg: "success",
	}

	return
}

func (a *Authorize) KeepLine(ctx context.Context, req *proto.Request) (resp *proto.Response, err error) {
	// 检查应用是否授权到期
	if !LoadLic().CheckTime(req.App) {
		resp = &proto.Response{Code: statusExpiresErr, Msg: "app does not exist or authorization expires"}
		return
	}
	// 检查实例是否存在
	cli, exist := GetClient(req.App + "/" + req.Id)
	if cli == nil || !exist {
		resp = &proto.Response{Code: statusGetClientErr, Msg: "the client does not exist or get error"}
		return
	}
	if utils.Md5sum([]byte(cli.Uuid), nil) != req.Umd5 || cli.Lease != req.Lease {
		resp = &proto.Response{Code: statusUmd5Err, Msg: "uuid md5sum error"}
		return
	}
	if err = KeepAliveClient(req.App+"/"+req.Id, req.Lease); err != nil {
		resp = &proto.Response{Code: statusKeepErr, Msg: "keep line error"}
		return
	}
	resp = &proto.Response{
		Code: statusOk,
		Data: &proto.Data{
			Auth:   nil,
			Cipher: nil,
			Lease:  req.Lease,
		},
		Msg: "success",
	}
	return
}

func (a *Authorize) OffLine(ctx context.Context, req *proto.Request) (resp *proto.Response, err error) {
	// 检查应用是否授权到期
	if !LoadLic().CheckTime(req.App) {
		resp = &proto.Response{Code: statusExpiresErr, Msg: "app does not exist or authorization expires"}
		return
	}
	// 检查实例是否存在
	cli, exist := GetClient(req.App + "/" + req.Id)
	if cli == nil || !exist {
		resp = &proto.Response{Code: statusGetClientErr, Msg: "the client does not exist"}
		return
	}
	if utils.Md5sum([]byte(cli.Uuid), nil) != req.Umd5 || cli.Lease != req.Lease {
		resp = &proto.Response{Code: statusUmd5Err, Msg: "uuid md5sum error"}
		return
	}
	if err = DelClient(req.App+"/"+req.Id, req.Lease); err != nil {
		resp = &proto.Response{Code: statusOffErr, Msg: "off line error"}
		return
	}
	resp = &proto.Response{
		Code: statusOk,
		Data: &proto.Data{
			Auth:   nil,
			Cipher: nil,
			Lease:  0,
		},
		Msg: "success",
	}
	return
}

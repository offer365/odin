package logic

import (
	"context"
	"github.com/offer365/odin/dao"
	"github.com/offer365/odin/embedder"
	"github.com/offer365/odin/model"
	"github.com/offer365/odin/node"
	"time"
)

const (
	licenseKey            = "/odin/license"
	clearLicenseKey       = "/odin/clear_license"
	clientConfigKeyPrefix = "/odin/client_config/"
	clientKeyPrefix       = "/odin/client/"
	serialNumKey          = "/odin/serial_num"
	membersKey            = "members"
	defaultKey            = "default"
)

var (
	store         dao.Store
	Device        embedder.Embed
	Serial        *model.SerialNum
	Self          *node.Node
	members       = make(map[string]string, 0)
	confWhiteList = map[string]string{"default": "", "members": ""} // 在白名单的配置无法删除
	PutWhiteList  = map[string]string{"members": ""}                // 在白名单的配置外部无法编辑
)

func init() {
	store = dao.NewStore()
	Serial = new(model.SerialNum)

}

func InitNode(addr, group, rpc string, peers []string) {
	Self = node.NewNode(addr, group, rpc, peers)
}

func InitStore(ip, port, user, pwd string, timeout time.Duration) (err error) {
	return store.Init(context.Background(), dao.WithHost(ip), dao.WithPort(port), dao.WithUsername(user), dao.WithPassword(pwd), dao.WithTimeout(timeout))
}

func InitEmbed(group, dir, addr, client, peer, state, metrics string, peers []string) (err error) {
	Device = embedder.NewEmbed()
	return Device.Init(context.Background(),
		embedder.WithGroup(group),
		embedder.WithDir(dir),
		embedder.WithIP(addr),
		embedder.WithClientPort(client),
		embedder.WithPeerPort(peer),
		embedder.WithCluster(peers),
		embedder.WithClusterState(state),
		embedder.WithMetrics(metrics, "b"))

}

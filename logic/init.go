package logic

import (
	"context"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/dao"
	"github.com/offer365/odin/embedder"
	"github.com/offer365/odin/model"
	"github.com/offer365/odin/node"
	"time"
)

const (
	licenseKey      = "/odin/license"
	clearLicenseKey = "/odin/clear_license"

	clientConfigKeyPrefix = "/odin/client_config/"
	clientKeyPrefix       = "/odin/client/"

	rankCodeKey = "/odin/rank_code"
)

var (
	store  dao.Store
	Device embedder.Embed
	//license *model.License
	Serial *model.SerialNum
	Self   *node.Node
)

func init() {
	store = dao.NewStore()
	//License = new(model.License)
	Serial = new(model.SerialNum)
	Self = node.NewNode(config.Cfg.Name, config.Cfg.Addr)
}

func InitStore(ip, port, user, pwd string, timeout time.Duration) (err error) {
	return store.Init(context.Background(), dao.WithHost(ip), dao.WithPort(port), dao.WithUsername(user), dao.WithPassword(pwd), dao.WithTimeout(timeout))
}

func InitEmbed() (err error) {
	Device = embedder.NewEmbed()
	return Device.Init(context.Background(),
		embedder.WithName(config.Cfg.Name),
		embedder.WithDir(config.Cfg.Dir),
		embedder.WithIP(config.Cfg.Addr),
		embedder.WithClientPort(config.Cfg.Client),
		embedder.WithPeerPort(config.Cfg.Peer),
		embedder.WithCluster(config.Cfg.Peers),
		embedder.WithClusterState(config.Cfg.State))
}

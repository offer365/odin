package logic

import (
	"context"
	"time"

	"github.com/offer365/example/etcd/dao"
	"github.com/offer365/example/etcd/embedder"
	"github.com/offer365/odin/log"
)

const (
	defaultKey            = "default"
)

var (
	store         dao.Store
	Device        embedder.Embed
	confWhiteList = map[string]string{"default": "", "members": ""} // 在白名单的配置无法删除
)

func InitStore(addr, user, pwd string, timeout time.Duration) (err error) {
	store = dao.NewStore()
	return store.Init(context.Background(),
		dao.WithAddr(addr),
		dao.WithUsername(user),
		dao.WithPassword(pwd),
		dao.WithTimeout(timeout),
	)
}

func InitEmbed(name, dir, client, peer, token, state string, cluster map[string]string) (err error) {
	Device = embedder.NewEmbed()
	return Device.Init(context.Background(),
		embedder.WithName(name),
		embedder.WithDir(dir),
		embedder.WithClientAddr(client),
		embedder.WithPeerAddr(peer),
		embedder.WithClusterToken(token),
		embedder.WithClusterState(state),
		embedder.WithCluster(cluster),
		embedder.WithLogger(log.Sugar),
	)
}

func Close() {
	Device.Close()
	store.Close()
}

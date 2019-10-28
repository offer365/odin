package logic

import (
	"context"
	"github.com/offer365/example/etcd/dao"
	"github.com/offer365/example/etcd/embedder"
	"time"
)

const (
	licenseKey            = "/M0n4c7gcVbw/TEzI9ZLhTJZ9AVB"
	clearLicenseKey       = "/RzgECgYEA6bORkz/jEW4KBteRzPi011"
	clientConfigKeyPrefix = "/GMyEfmKOR4vYGq73/XWNjx5XXmoVDhU/"
	clientKeyPrefix       = "/MIIEvwIBADAN/JtfE6TPrB8WFpEdgJKRir/"
	tokenKey              = "/xQs3d1n7Ye81v/NBT2he8EJWiBzyZ31/"
	serialNumKey          = "/hIrNugCu9vQC8PYDkfWelXPJxI/e6gIs8DJx9WF02JD13Yi"
	defaultKey            = "default"
)

var (
	store  dao.Store
	Device embedder.Embed
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

func InitEmbed(name, dir, client, peer, state string, cluster map[string]string) (err error) {
	Device = embedder.NewEmbed()
	return Device.Init(context.Background(),
		embedder.WithName(name),
		embedder.WithDir(dir),
		embedder.WithClientAddr(client),
		embedder.WithPeerAddr(peer),
		embedder.WithClusterState(state),
		embedder.WithCluster(cluster),
	)
}

func Close() {
	Device.Close()
	store.Close()
}

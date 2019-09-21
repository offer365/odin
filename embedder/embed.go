package embedder

import (
	"context"
	"strings"
)

type Options struct {
	Name         string
	Dir          string
	IP           string
	ClientPort   string
	PeerPort     string
	Cluster      []string
	ClusterState string // "new" or "existing"
}

type Option func(opts *Options)

func WithName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func WithDir(dir string) Option {
	return func(opts *Options) {
		opts.Dir = dir
	}
}

func WithIP(ip string) Option {
	return func(opts *Options) {
		opts.IP = ip
	}
}

func WithClientPort(cp string) Option {
	return func(opts *Options) {
		opts.ClientPort = cp
	}
}

func WithPeerPort(pp string) Option {
	return func(opts *Options) {
		opts.PeerPort = pp
	}
}

func WithCluster(clu []string) Option {
	return func(opts *Options) {
		opts.Cluster = clu
	}
}

func WithClusterState(state string) Option {
	return func(opts *Options) {
		// "new" or "existing"
		if strings.HasPrefix(state, "exist") {
			opts.ClusterState = "existing"
		} else {
			opts.ClusterState = "new"
		}
	}
}

type Embed interface {
	Init(ctx context.Context, option ...Option) (err error)
	Run(ready chan struct{}) (err error)
	SetAuth(username, password string) (err error)
	IsLeader() bool
}

func NewEmbed() Embed {
	return new(etcdEmbed)
}

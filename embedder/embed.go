package embedder

import (
	"context"
	"strings"
)

type Options struct {
	name         string
	dir          string
	ip           string
	clientPort   string
	peerPort     string
	cluster      []string
	clusterState string // "new" or "existing"
	metrics      string
	metricsUrl   string
}

type Option func(opts *Options)

func WithName(name string) Option {
	return func(opts *Options) {
		opts.name = name
	}
}

func WithDir(dir string) Option {
	return func(opts *Options) {
		opts.dir = dir
	}
}

func WithIP(ip string) Option {
	return func(opts *Options) {
		opts.ip = ip
	}
}

func WithClientPort(cp string) Option {
	return func(opts *Options) {
		opts.clientPort = cp
	}
}

func WithPeerPort(pp string) Option {
	return func(opts *Options) {
		opts.peerPort = pp
	}
}

func WithCluster(clu []string) Option {
	return func(opts *Options) {
		opts.cluster = clu
	}
}

func WithClusterState(state string) Option {
	return func(opts *Options) {
		// "new" or "existing"
		if strings.HasPrefix(state, "exist") {
			opts.clusterState = "existing"
		} else {
			opts.clusterState = "new"
		}
	}
}

func WithMetrics(addr string, mode string) Option {
	return func(opts *Options) {
		switch {
		case strings.HasPrefix(mode, "b"):
			opts.metrics = "base"
		case strings.HasPrefix(mode, "e"):
			opts.metrics = "extensive"
		default:
			opts.metrics = "base"
		}
		if addr != "" && !strings.HasPrefix(addr, "http://") {
			opts.metricsUrl = "http://" +opts.ip+":"+ addr
			return
		}
		opts.metricsUrl = addr
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

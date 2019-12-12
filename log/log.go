package log

import (
	"go.etcd.io/etcd/pkg/logutil"
	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func init() {
	var (
		lg  *zap.Logger
		err error
	)
	if lg, err = zap.NewProduction(); err != nil {
		return
	}
	defer lg.Sync()
	cfg := logutil.DefaultZapLoggerConfig
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	if lg, err = cfg.Build(); err != nil {
		return
	}
	Sugar = lg.Sugar()
}

package log

import (
	"go.etcd.io/etcd/pkg/logutil"
	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func init() {
	lg, _ := zap.NewProduction()
	defer lg.Sync()
	cfg := logutil.DefaultZapLoggerConfig
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	lg, _ = cfg.Build()
	Sugar = lg.Sugar()
}

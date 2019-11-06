package log

import (
	"github.com/offer365/example/loger"
	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func init() {
	Sugar, _ = loger.SugaredLog("ODIN_RUN_MODE", loger.ReleaseMode)
}

package log

import "testing"

func TestLog(t *testing.T) {
	log.Sugar.Info("test")
	log.Sugar.Warn(map[string]string{"a": "b"})
}

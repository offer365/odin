package log

import "testing"

func TestLog(t *testing.T) {
	Sugar.Info("test")
	Sugar.Warn(map[string]string{"a": "b"})
}

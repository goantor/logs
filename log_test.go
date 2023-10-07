package logs

import (
	"testing"
)

func TestLog(t *testing.T) {
	opt := &Options{
		Path:            "./",
		Level:           "debug",
		Stdout:          true,
		SaveDay:         1,
		TimestampFormat: "2006-01-02T15:04:06",
	}

	en := NewEntity(opt)

	en.Initialize()

	l := New("/test", "POST", "192.168.1.1")

	l.Error("xxxx", map[string]interface{}{
		"name": "ll",
	})
}

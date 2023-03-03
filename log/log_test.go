package log

import (
	"os"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestMain(t *testing.M) {
	writer := &lumberjack.Logger{
		Filename:   "tmp.log",
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 30,
		LocalTime:  true,
		Compress:   false,
	}
	Init(writer)
	os.Exit(t.Run())
}

func TestLog(t *testing.T) {
}

package main

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jdxj/oh-my-feed/log"
)

var (
	logger *zap.SugaredLogger
)

func main() {
	conf := &lumberjack.Logger{
		Filename:   "/Users/jdxj/workspace/oh-my-feed/cmd/test/tmp.log",
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 30,
		LocalTime:  true,
		Compress:   false,
	}
	log.Init(conf)
	logger = log.Logger

	logger.Infof("hh")
}

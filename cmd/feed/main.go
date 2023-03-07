package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jdxj/oh-my-feed/internal/app/bot"
	"github.com/jdxj/oh-my-feed/internal/app/model"
	"github.com/jdxj/oh-my-feed/internal/app/task"
	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

var (
	configPath = flag.String("config", "config.yaml", "config path")
)

func main() {
	flag.Parse()

	config.Init(*configPath)
	log.Init()
	model.Init()
	task.Init()
	bot.Init()
	log.Infof("started")

	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)
	<-signs

	log.Infof("stopped")
	bot.Stop()
	task.Stop()
	log.Sync()
}

package task

import (
	"os"
	"testing"
	"time"

	"github.com/jdxj/oh-my-feed/internal/app/model"
	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

func TestMain(t *testing.M) {
	config.Init("../config/config.yaml")
	log.Init()
	model.Init()

	os.Exit(t.Run())
}

func TestInit(t *testing.T) {
	Init()

	time.Sleep(time.Second * 30)

	Stop()
}

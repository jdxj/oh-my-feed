package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jdxj/oh-my-feed/internal/pkg/config"
)

func TestMain(t *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(wd, "../config/config.yaml")
	config.Init(configPath)
	Init()
	os.Exit(t.Run())
}

func TestLog(t *testing.T) {
	Debugf("%s", "abc")
}

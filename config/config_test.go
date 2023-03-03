package config

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	Init("config.yaml")
	os.Exit(t.Run())
}

func TestInit(t *testing.T) {
	fmt.Printf("%+v\n%+v\n%+v\n", DB, Telegram, Logger)
}

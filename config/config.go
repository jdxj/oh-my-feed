package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type db struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}

type telegram struct {
	Token string `yaml:"token"`
}

type logger struct {
	Filename string `yaml:"filename"`
	MaxAge   int    `yaml:"maxage"`
	Level    int    `yaml:"level"`
}

var (
	DB       = &db{}
	Telegram = &telegram{}
	Logger   = &logger{}

	global = &struct {
		DB       *db       `yaml:"db"`
		Telegram *telegram `yaml:"telegram"`
		Logger   *logger   `yaml:"logger"`
	}{
		DB:       DB,
		Telegram: Telegram,
		Logger:   Logger,
	}
)

func Init(name string) {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(global)
	if err != nil {
		panic(err)
	}
}

package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config interface {
	Init()
}

type Configs []Config

func InitAll(configs Configs) {
	InitConfigFile()
	for i := range configs {
		configs[i].Init()
	}
}

func InitConfigFile() {
	viper.SetConfigFile(os.Getenv("CONFIG_PATH"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

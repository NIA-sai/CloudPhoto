package config

import (
	"CloudPhoto/internal/tool"
	"github.com/spf13/viper"
)

var c config

func init() {
	viper.SetConfigFile("./config.yaml")
}
func Read() {
	tool.PanicIfErr(
		viper.ReadInConfig(),
		viper.Unmarshal(&c),
	)
}

func Get() config {
	return c
}

type config struct {
	App struct {
		Host string
		Mode string
		Port int
	}

	Mysql struct {
	}

	Redis struct {
	}

	Ai struct {
	}
}

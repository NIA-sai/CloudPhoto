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
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Ai struct {
		CutOut volcengineApi `mapstructure:"cut-out"`
	} `mapstructure:"ai"`
}

type volcengineApi struct {
	Url         string `mapstructure:"url"`
	Region      string `mapstructure:"region"`
	Service     string `mapstructure:"service"`
	Method      string `mapstructure:"method"`
	ContentType string `mapstructure:"content-type"`
	Action      string `mapstructure:"action"`
	Version     string `mapstructure:"version"`
	AccessId    string `mapstructure:"access-id"`
	SecretKey   string `mapstructure:"secret-key"`
}

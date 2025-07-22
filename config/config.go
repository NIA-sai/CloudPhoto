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
		CutOut     volcengineApi   `mapstructure:"cut-out"`
		FaceFusion tencentCloudApi `mapstructure:"face-fusion"`
	} `mapstructure:"ai"`
}

type tencentCloudApi struct { //配置结构体
	Url       string `mapstructure:"url"`        // API地址
	SecretId  string `mapstructure:"secret-id"`  // 密钥ID
	SecretKey string `mapstructure:"secret-key"` // 密钥Key
	ProjectId string `mapstructure:"project-id"` // 活动ID（必填）
	Action    string `mapstructure:"action"`     // 接口动作（必填，固定值）
	Version   string `mapstructure:"version"`    // 接口版本（必填，固定值）
	Region    string `mapstructure:"region"`     // 地域（必填）
	//RspImgType string `mapstructure:"rsp-img-type"` // 返回图像格式（必填）
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

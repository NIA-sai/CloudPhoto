package storage

import (
	"CloudPhoto/config"
	"github.com/mojocn/base64Captcha"
)

func GetBodyFilePath(id string) string {
	return "./body"
}

func Init() {
	if config.Get().App.CaptchaUseTimes > 1 {
		captcha = &MultiUseStore{}
	} else {
		captcha = base64Captcha.DefaultMemStore
	}
}

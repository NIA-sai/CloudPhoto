package storage

import (
	"CloudPhoto/config"
	"github.com/mojocn/base64Captcha"
	"time"
)

func GetBodyFilePath(id string) string {
	return "./body"
}

func Init() {
	useTimes := config.Get().App.CaptchaUseTimes
	if useTimes > 1 {
		captcha = NewMultipleUseCaptcha(1024, 10*time.Minute, useTimes)
	} else {
		captcha = base64Captcha.DefaultMemStore
	}
}

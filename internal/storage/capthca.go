package storage

import "github.com/mojocn/base64Captcha"

var captcha = base64Captcha.DefaultMemStore

func GetCaptcha() *base64Captcha.Store {

	return &captcha
}

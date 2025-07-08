package captcha

import (
	"CloudPhoto/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Captcha struct {
}

func (*Captcha) GetName() string {
	return "captcha"
}
func (*Captcha) Init() {
}

func (*Captcha) InitRouter(r *gin.RouterGroup) {
	r.Use(middleware.CaptchaAuth())
	r.GET("/image", ImageCaptcha)
}

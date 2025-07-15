package captcha

import (
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
	r.GET("/image", imageCaptcha)
}

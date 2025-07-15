package aitool

import (
	"CloudPhoto/config"
	//"CloudPhoto/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/url"
)

type AiTool struct {
}

func (*AiTool) GetName() string {
	return "api"
}
func (*AiTool) Init() {
	rawCutOutUrl, _ = url.Parse(config.Get().Ai.CutOut.Url)
	query := url.Values{}
	query.Add("Action", config.Get().Ai.CutOut.Action)
	query.Add("Version", config.Get().Ai.CutOut.Version)
	rawCutOutUrl.RawQuery = query.Encode()
	cutOutUrl = rawCutOutUrl.String()
}

func (*AiTool) InitRouter(r *gin.RouterGroup) {
	//r.Use(middleware.CaptchaAuth())
	r.POST("/changeFace", changeFace)
	r.POST("/cutOutFigure", cutOutFigure)
}

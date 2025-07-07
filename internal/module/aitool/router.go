package aitool

import "github.com/gin-gonic/gin"

type AiTool struct {
}

func (*AiTool) GetName() string {
	return "api"
}
func (*AiTool) Init() {
}

func (*AiTool) InitRouter(r *gin.RouterGroup) {
	r.POST("/changeFace", ChangeFace)
	r.POST("/addFigure", AddFigure)
}

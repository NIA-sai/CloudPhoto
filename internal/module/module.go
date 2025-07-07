package module

import (
	"CloudPhoto/internal/module/aitool"
	"github.com/gin-gonic/gin"
)

type Module interface {
	GetName() string
	Init()
	InitRouter(r *gin.RouterGroup)
}

var Modules []Module

func registerModule(m ...Module) {
	Modules = append(Modules, m...)
}

func init() {
	registerModule(
		&aitool.AiTool{},
	)
}

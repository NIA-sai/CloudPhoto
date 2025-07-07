package server

import (
	"CloudPhoto/config"
	"github.com/gin-gonic/gin"
)

func Init() {
	config.Read()
}
func Start() {
	gin.SetMode(config.Mode())
	r := gin.Default()
}
func Stop() {

}

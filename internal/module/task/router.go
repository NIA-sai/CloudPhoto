package task

import "github.com/gin-gonic/gin"

type Task struct {
}

func (*Task) GetName() string {
	return "task"
}
func (*Task) Init() {
}

func (*Task) InitRouter(r *gin.RouterGroup) {
	r.GET("/ask/:type/:taskId", taskAsk)
}

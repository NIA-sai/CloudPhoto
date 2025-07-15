package task

import (
	"CloudPhoto/internal/storage"
	"github.com/gin-gonic/gin"
)

func taskAsk(c *gin.Context) {
	typo := c.Param("type")
	taskId := c.Param("taskId")
	var result string
	switch typo {
	case storage.ChangeFace:
		//result, ok = storage.GetChangeFaceTask(taskId)
	case storage.CutOut:
		result = storage.GetCutOutTask(taskId)
	default:
		c.Status(404)
		return
	}
	c.Data(200, "application/json", []byte(result))
	return
}

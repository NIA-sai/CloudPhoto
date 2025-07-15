package task

import (
	"CloudPhoto/internal/model"
	"CloudPhoto/internal/storage"
	"github.com/gin-gonic/gin"
)

func taskAsk(c *gin.Context) {
	typo := c.Param("type")
	taskId := c.Param("taskId")
	var result model.TaskResult
	var ok bool
	switch typo {
	case storage.ChangeFace:
		result, ok = storage.GetChangeFaceTask(taskId)
	case storage.CutOut:
		result, ok = storage.GetCutOutTask(taskId)
	default:
		c.Status(404)
		return
	}
	if !ok {
		c.Status(404)
		return
	}
	c.JSON(200, result)
	return
}

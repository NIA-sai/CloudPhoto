package storage

import (
	"CloudPhoto/internal/model"
	"sync"
)

// 也许可以改成使用redis?
type taskStore = struct {
	sync.Mutex
	data map[string]model.TaskResult
}

var (
	changeFaceTaskStore = taskStore{
		data: make(map[string]model.TaskResult),
	}
	cutOutTaskStore = taskStore{
		data: make(map[string]model.TaskResult),
	}
)

func SetChangeFaceTask(taskId string, result model.TaskResult) {
	setTask(&changeFaceTaskStore, taskId, result)
}
func GetChangeFaceTask(id string) (model.TaskResult, bool) {
	return getTask(&changeFaceTaskStore, id)
}

func SetCutOutTask(taskId string, result model.TaskResult) {
	setTask(&cutOutTaskStore, taskId, result)
}
func GetCutOutTask(id string) (model.TaskResult, bool) {
	return getTask(&cutOutTaskStore, id)
}

func setTask(ts *taskStore, taskId string, result model.TaskResult) {
	ts.Lock()
	defer ts.Unlock()
	ts.data[taskId] = result
}

func getTask(ts *taskStore, id string) (model.TaskResult, bool) {
	ts.Lock()
	defer ts.Unlock()
	val, ok := ts.data[id]
	return val, ok
}

package storage

import (
	"CloudPhoto/internal/database/redis"
	"CloudPhoto/internal/model"
	"context"
	"encoding/json"
	"sync"
	"time"
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
	//存内存
	//setTask(&cutOutTaskStore, taskId, result)
	//存redis
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	redis.Set(&ctx, "cp:cc:"+taskId, data, 24*time.Hour)
}
func GetCutOutTask(id string) string {
	//return getTask(&cutOutTaskStore, id)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return redis.Get(&ctx, "cp:cc:"+id)
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

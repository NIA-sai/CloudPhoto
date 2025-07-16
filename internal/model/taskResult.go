package model

type TaskResult struct {
	Status  string      `json:"status"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
}

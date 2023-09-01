package models

import "time"

type TaskPage struct {
	Tasks []Task
	Task  string
	Error string
	Key   int64
}

func NewTaskPage(tasks []Task) TaskPage {
	return TaskPage{Key: time.Now().UnixMilli(), Tasks: tasks}
}

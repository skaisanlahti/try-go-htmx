package models

type Task struct {
	Id   int
	Task string
	Done bool
}

func NewTask(task string) Task {
	return Task{Task: task, Done: false}
}

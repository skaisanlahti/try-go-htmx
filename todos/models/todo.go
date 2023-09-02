package models

type Todo struct {
	Id   int
	Task string
	Done bool
}

func NewTask(task string) Todo {
	return Todo{Task: task, Done: false}
}

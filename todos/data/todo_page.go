package data

import "time"

type TodoPage struct {
	Todos []Todo
	Task  string
	Error string
	Key   int64
}

func NewTodoPage() TodoPage {
	return TodoPage{Key: time.Now().UnixMilli()}
}

package entity

import "errors"

type Todo struct {
	Id   int
	Task string
	Done bool
}

func NewTodo(task string) Todo {
	return Todo{Task: task}
}

func (this Todo) Toggle() Todo {
	this.Done = !this.Done
	return this
}

func (this Todo) Validate() error {
	length := len([]rune(this.Task))
	if length == 0 {
		return errors.New("Task is too short.")
	}

	if length > 100 {
		return errors.New("Task is too long.")
	}

	return nil
}

package entity

import "errors"

type Todo struct {
	Id   int
	Task string
	Done bool
}

func NewTodo(task string) (Todo, error) {
	length := len([]rune(task))
	if length == 0 {
		return Todo{}, errors.New("Task is too short.")
	}

	if length > 100 {
		return Todo{}, errors.New("Task is too long.")
	}

	return Todo{Task: task}, nil
}

package domain

import "errors"

type Todo struct {
	Id   int
	Task string
	Done bool
}

func NewTodo(task string) (Todo, error) {
	if task == "" {
		return Todo{}, errors.New("Task too short.")
	}

	if len(task) > 100 {
		return Todo{}, errors.New("Task too long.")
	}

	return Todo{Task: task, Done: false}, nil
}

func ToggleTodo(todo Todo) Todo {
	todo.Done = !todo.Done
	return todo
}

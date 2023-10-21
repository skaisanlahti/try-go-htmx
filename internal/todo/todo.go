package todo

import "errors"

type todo struct {
	Id   int
	Task string
	Done bool
}

func newTodo(task string) (todo, error) {
	length := len([]rune(task))
	if length == 0 {
		return todo{}, errors.New("Task is too short.")
	}

	if length > 100 {
		return todo{}, errors.New("Task is too long.")
	}

	return todo{Task: task}, nil
}

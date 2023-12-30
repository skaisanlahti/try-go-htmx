package entity

import "errors"

type Todo struct {
	Id         int
	Task       string
	Done       bool
	TodoListId int
}

func NewTodo(task string, todoListId int) Todo {
	return Todo{Task: task, TodoListId: todoListId}
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

type TodoList struct {
	Id     int
	Name   string
	UserId int
}

func NewTodoList(name string, userId int) TodoList {
	return TodoList{Name: name, UserId: userId}
}

func (this TodoList) Validate() error {
	length := len([]rune(this.Name))
	if length == 0 {
		return errors.New("Name is too short.")
	}

	if length > 100 {
		return errors.New("Name is too long.")
	}

	return nil
}

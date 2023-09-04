package interfaces

import "github.com/skaisanlahti/try-go-htmx/todos/data"

type DataService interface {
	FindTodos() ([]data.Todo, error)
	FindTodoByID(id int) (data.Todo, error)
	AddTodo(todo data.Todo) error
	UpdateTodo(todo data.Todo) error
	RemoveTodo(id int) error
}

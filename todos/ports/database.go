package ports

import "github.com/skaisanlahti/try-go-htmx/todos/models"

type Database interface {
	GetTodos() ([]models.Todo, error)
	GetTodoByID(id int) (models.Todo, error)
	AddTodo(todo models.Todo) error
	UpdateTodo(todo models.Todo) error
	RemoveTodo(id int) error
}

package todo

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type TodoService struct {
	todos *todoStorage
}

func NewTodoService(database *sql.DB) *TodoService {
	return &TodoService{newTodoStorage(database)}
}

func (service *TodoService) FindTodos() []entity.Todo {
	return service.todos.findTodos()
}

func (service *TodoService) FindTodo(id int) (entity.Todo, error) {
	return service.todos.findTodoById(id)
}

func (service *TodoService) AddTodo(task string) error {
	newTodo, err := entity.NewTodo(task)
	if err != nil {
		return err
	}

	if err := service.todos.insertTodo(newTodo); err != nil {
		return err
	}

	return nil
}

func (service *TodoService) ToggleTodo(todoId int) error {
	todo, err := service.todos.findTodoById(todoId)
	if err != nil {
		return err
	}

	todo.Done = !todo.Done
	if err := service.todos.updateTodo(todo); err != nil {
		return err
	}

	return nil
}

func (service *TodoService) RemoveTodo(todoId int) error {
	_, err := service.todos.findTodoById(todoId)
	if err != nil {
		return err
	}

	if err := service.todos.deleteTodo(todoId); err != nil {
		return err
	}

	return nil
}

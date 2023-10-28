package todo

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type TodoService struct {
	storage *todoStorage
}

func NewTodoService(database *sql.DB) *TodoService {
	return &TodoService{newTodoStorage(database)}
}

func (this *TodoService) FindTodos() []entity.Todo {
	return this.storage.findTodos()
}

func (this *TodoService) FindTodoById(id int) (entity.Todo, error) {
	return this.storage.findTodoById(id)
}

func (this *TodoService) AddTodo(task string) error {
	newTodo := entity.NewTodo(task)
	err := newTodo.Validate()
	if err != nil {
		return err
	}

	if err := this.storage.insertTodo(newTodo); err != nil {
		return err
	}

	return nil
}

func (this *TodoService) ToggleTodo(todoId int) error {
	todo, err := this.storage.findTodoById(todoId)
	if err != nil {
		return err
	}

	if err := this.storage.updateTodo(todo.Toggle()); err != nil {
		return err
	}

	return nil
}

func (this *TodoService) RemoveTodo(todoId int) error {
	_, err := this.storage.findTodoById(todoId)
	if err != nil {
		return err
	}

	if err := this.storage.deleteTodo(todoId); err != nil {
		return err
	}

	return nil
}

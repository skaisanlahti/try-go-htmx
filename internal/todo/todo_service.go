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

func (this *TodoService) FindListById(listId int) (entity.TodoList, error) {
	return this.storage.findTodoListById(listId)
}

func (this *TodoService) FindListsByUserId(userId int) []entity.TodoList {
	return this.storage.findTodoListsByUserId(userId)
}

func (this *TodoService) FindTodosByListId(listId int) []entity.Todo {
	return this.storage.findTodosByListId(listId)
}

func (this *TodoService) FindTodoById(id int) (entity.Todo, error) {
	return this.storage.findTodoById(id)
}

func (this *TodoService) AddList(name string, userId int) (entity.TodoList, error) {
	newList := entity.NewTodoList(name, userId)
	if err := newList.Validate(); err != nil {
		return newList, err
	}

	if err := this.storage.insertTodoList(newList); err != nil {
		return newList, err
	}

	return newList, nil
}

func (this *TodoService) RemoveList(listId int) (entity.TodoList, error) {
	list, err := this.storage.findTodoListById(listId)
	if err != nil {
		return list, err
	}

	if err := this.storage.deleteTodoList(listId); err != nil {
		return list, err
	}

	return list, nil
}

func (this *TodoService) AddTodo(task string, todoListId int) (entity.Todo, error) {
	newTodo := entity.NewTodo(task, todoListId)
	if err := newTodo.Validate(); err != nil {
		return newTodo, err
	}

	if err := this.storage.insertTodo(newTodo); err != nil {
		return newTodo, err
	}

	return newTodo, nil
}

func (this *TodoService) ToggleTodo(todoId int) (entity.Todo, error) {
	todo, err := this.storage.findTodoById(todoId)
	if err != nil {
		return todo, err
	}

	newTodo := todo.Toggle()
	if err := this.storage.updateTodo(newTodo); err != nil {
		return todo, err
	}

	return newTodo, nil
}

func (this *TodoService) RemoveTodo(todoId int) (entity.Todo, error) {
	todo, err := this.storage.findTodoById(todoId)
	if err != nil {
		return todo, err
	}

	if err := this.storage.deleteTodo(todoId); err != nil {
		return todo, err
	}

	return todo, nil
}

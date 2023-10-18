package todoapp

import "errors"

type Todo struct {
	Id   int
	Task string
	Done bool
}

func CreateTodo(task string) (Todo, error) {
	length := len([]rune(task))
	if length == 0 {
		return Todo{}, errors.New("Task is too short.")
	}

	if length > 100 {
		return Todo{}, errors.New("Task is too long.")
	}

	return Todo{Task: task}, nil
}

type TodoStorage interface {
	FindTodos() []Todo
	FindTodoById(id int) (Todo, error)
	AddTodo(todo Todo) error
	UpdateTodo(todo Todo) error
	RemoveTodo(id int) error
}

type TodoService struct {
	Settings    Settings
	TodoStorage TodoStorage
}

func CreateTodoService(s Settings, t TodoStorage) *TodoService {
	return &TodoService{s, t}
}

func (service *TodoService) AddTodo(task string) error {
	newTodo, err := CreateTodo(task)
	if err != nil {
		return err
	}

	if err := service.TodoStorage.AddTodo(newTodo); err != nil {
		return err
	}

	return nil
}

func (service *TodoService) ToggleTodo(todoId int) error {
	todo, err := service.TodoStorage.FindTodoById(todoId)
	if err != nil {
		return err
	}

	todo.Done = !todo.Done
	if err := service.TodoStorage.UpdateTodo(todo); err != nil {
		return err
	}

	return nil
}

func (service *TodoService) RemoveTodo(todoId int) error {
	_, err := service.TodoStorage.FindTodoById(todoId)
	if err != nil {
		return err
	}

	if err := service.TodoStorage.RemoveTodo(todoId); err != nil {
		return err
	}

	return nil
}

package todoapp

import "errors"

var (
	ErrTodoTaskTooShort = errors.New("Task is too short.")
	ErrTodoTaskTooLong  = errors.New("Task is too long.")
)

type Todo struct {
	Id   int
	Task string
	Done bool
}

func NewTodo(task string) (Todo, error) {
	length := len([]rune(task))
	if length == 0 {
		return Todo{}, ErrTodoTaskTooShort
	}

	if length > 100 {
		return Todo{}, ErrTodoTaskTooLong
	}

	return Todo{Task: task}, nil
}

type TodoAccessor interface {
	FindTodos() []Todo
	FindTodoById(id int) (Todo, error)
	AddTodo(todo Todo) error
	UpdateTodo(todo Todo) error
	RemoveTodo(id int) error
}

type TodoWriter struct {
	Settings     Settings
	TodoAccessor TodoAccessor
}

func NewTodoWriter(s Settings, t TodoAccessor) *TodoWriter {
	return &TodoWriter{s, t}
}

func (writer *TodoWriter) AddTodo(task string) error {
	newTodo, err := NewTodo(task)
	if err != nil {
		return err
	}

	if err := writer.TodoAccessor.AddTodo(newTodo); err != nil {
		return err
	}

	return nil
}

func (writer *TodoWriter) ToggleTodo(todoId int) error {
	todo, err := writer.TodoAccessor.FindTodoById(todoId)
	if err != nil {
		return err
	}

	todo.Done = !todo.Done
	if err := writer.TodoAccessor.UpdateTodo(todo); err != nil {
		return err
	}

	return nil
}

func (writer *TodoWriter) RemoveTodo(todoId int) error {
	_, err := writer.TodoAccessor.FindTodoById(todoId)
	if err != nil {
		return err
	}

	if err := writer.TodoAccessor.RemoveTodo(todoId); err != nil {
		return err
	}

	return nil
}

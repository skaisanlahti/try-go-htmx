package todo

type todoService struct {
	todoStorage *todoStorage
}

func NewTodoService(t *todoStorage) *todoService {
	return &todoService{t}
}

func (service *todoService) addTodo(task string) error {
	newTodo, err := newTodo(task)
	if err != nil {
		return err
	}

	if err := service.todoStorage.addTodo(newTodo); err != nil {
		return err
	}

	return nil
}

func (service *todoService) toggleTodo(todoId int) error {
	todo, err := service.todoStorage.findTodoById(todoId)
	if err != nil {
		return err
	}

	todo.Done = !todo.Done
	if err := service.todoStorage.updateTodo(todo); err != nil {
		return err
	}

	return nil
}

func (service *todoService) removeTodo(todoId int) error {
	_, err := service.todoStorage.findTodoById(todoId)
	if err != nil {
		return err
	}

	if err := service.todoStorage.removeTodo(todoId); err != nil {
		return err
	}

	return nil
}

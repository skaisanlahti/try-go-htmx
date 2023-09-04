package services

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skaisanlahti/try-go-htmx/todos/data"
	"github.com/skaisanlahti/try-go-htmx/todos/interfaces"
)

type TodoHttpService struct {
	dataService     interfaces.DataService
	templateService interfaces.TemplateService
}

func NewTodoHttpService(
	dataService interfaces.DataService,
	templateService interfaces.TemplateService,
) *TodoHttpService {
	return &TodoHttpService{dataService, templateService}
}

func (this *TodoHttpService) GetTodoPage(response http.ResponseWriter, request *http.Request) error {
	todos, err := this.dataService.FindTodos()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.templateService.RenderPage(response, todos)
}

func (this *TodoHttpService) GetTodoList(response http.ResponseWriter, request *http.Request) error {
	todos, err := this.dataService.FindTodos()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.templateService.RenderList(response, todos)
}

func (this *TodoHttpService) AddTodo(response http.ResponseWriter, request *http.Request) error {
	task := request.FormValue("task")
	if task == "" {
		errorMessage := "Task can't be empty."
		this.templateService.RenderErrorForm(response, errorMessage)
		return nil
	}

	newTodo := data.NewTask(task)
	if err := this.dataService.AddTodo(newTodo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.Header().Add("HX-Trigger", "GetTodoList")
	return this.templateService.RenderEmptyForm(response)
}

func extractTodoID(url *url.URL) (int, error) {
	values := url.Query()
	idStr := values.Get("id")
	if idStr == "" {
		return 0, errors.New("Task ID not found in query")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (this *TodoHttpService) ToggleTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	todo, err := this.dataService.FindTodoByID(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	todo.Done = !todo.Done
	if err = this.dataService.UpdateTodo(todo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.templateService.RenderItem(response, todo)
}

func (this *TodoHttpService) RemoveTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = this.dataService.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.WriteHeader(http.StatusOK)
	return nil
}

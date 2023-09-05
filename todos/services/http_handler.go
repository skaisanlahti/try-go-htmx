package services

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skaisanlahti/try-go-htmx/todos/data"
)

type Storage interface {
	FindTodos() ([]data.Todo, error)
	FindTodoByID(id int) (data.Todo, error)
	AddTodo(todo data.Todo) error
	UpdateTodo(todo data.Todo) error
	RemoveTodo(id int) error
}

type Renderer interface {
	RenderPage(response http.ResponseWriter, todos []data.Todo) error
	RenderList(response http.ResponseWriter, todos []data.Todo) error
	RenderEmptyForm(response http.ResponseWriter) error
	RenderErrorForm(response http.ResponseWriter, errorMessage string) error
	RenderItem(response http.ResponseWriter, todo data.Todo) error
}

type HttpHandler struct {
	storage  Storage
	renderer Renderer
}

func NewHttpHandler(
	storage Storage,
	renderer Renderer,
) *HttpHandler {
	return &HttpHandler{storage, renderer}
}

func (this *HttpHandler) GetTodoPage(response http.ResponseWriter, request *http.Request) error {
	todos, err := this.storage.FindTodos()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.renderer.RenderPage(response, todos)
}

func (this *HttpHandler) GetTodoList(response http.ResponseWriter, request *http.Request) error {
	todos, err := this.storage.FindTodos()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.renderer.RenderList(response, todos)
}

func (this *HttpHandler) AddTodo(response http.ResponseWriter, request *http.Request) error {
	task := request.FormValue("task")
	if task == "" {
		errorMessage := "Task can't be empty."
		this.renderer.RenderErrorForm(response, errorMessage)
		return nil
	}

	newTodo := data.NewTask(task)
	if err := this.storage.AddTodo(newTodo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.Header().Add("HX-Trigger", "GetTodoList")
	return this.renderer.RenderEmptyForm(response)
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

func (this *HttpHandler) ToggleTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	todo, err := this.storage.FindTodoByID(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	todo.Done = !todo.Done
	if err = this.storage.UpdateTodo(todo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.renderer.RenderItem(response, todo)
}

func (this *HttpHandler) RemoveTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = this.storage.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.WriteHeader(http.StatusOK)
	return nil
}

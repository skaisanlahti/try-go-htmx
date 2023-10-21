package todo

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type htmxController struct {
	todoService  *todoService
	htmxRenderer *htmxRenderer
}

func NewHtmxController(service *todoService, renderer *htmxRenderer) *htmxController {
	return &htmxController{service, renderer}
}

func (controller *htmxController) getTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := controller.todoService.todoStorage.findTodos()
	html := controller.htmxRenderer.renderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (controller *htmxController) getTodoList(response http.ResponseWriter, request *http.Request) {
	todos := controller.todoService.todoStorage.findTodos()
	html := controller.htmxRenderer.renderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (controller *htmxController) addTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := controller.todoService.addTodo(task)
	if err != nil {
		html := controller.htmxRenderer.renderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	html := controller.htmxRenderer.renderTodoForm("", "")
	response.Header().Add("HX-Trigger", "GetTodoList")
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

var (
	ErrTodoIdMissing   = errors.New("Todo id not found in query.")
	ErrTodoIdNotNumber = errors.New("Todo id not a number.")
)

func extractTodoId(url *url.URL) (int, error) {
	values := url.Query()
	maybeId := values.Get("id")
	if maybeId == "" {
		return 0, ErrTodoIdMissing
	}

	id, err := strconv.Atoi(maybeId)
	if err != nil {
		log.Println(err.Error())
		return 0, ErrTodoIdNotNumber
	}

	return id, nil
}

func (controller *htmxController) toggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controller.todoService.toggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := controller.todoService.todoStorage.findTodoById(id)
	html := controller.htmxRenderer.renderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (controller *htmxController) removeTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controller.todoService.removeTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

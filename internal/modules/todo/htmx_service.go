package todo

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type htmxService struct {
	todoService  *todoService
	htmxRenderer *htmxRenderer
}

func newHtmxService(todoService *todoService, htmxRenderer *htmxRenderer) *htmxService {
	return &htmxService{todoService, htmxRenderer}
}

func (service *htmxService) getTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := service.todoService.todoStorage.findTodos()
	html := service.htmxRenderer.renderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) getTodoList(response http.ResponseWriter, request *http.Request) {
	todos := service.todoService.todoStorage.findTodos()
	html := service.htmxRenderer.renderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) addTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := service.todoService.addTodo(task)
	if err != nil {
		html := service.htmxRenderer.renderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	html := service.htmxRenderer.renderTodoForm("", "")
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

func (service *htmxService) toggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = service.todoService.toggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := service.todoService.todoStorage.findTodoById(id)
	html := service.htmxRenderer.renderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) removeTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = service.todoService.removeTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

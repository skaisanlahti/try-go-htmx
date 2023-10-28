package htmx

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (client *Client) getTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := client.app.FindTodos()
	html := client.renderer.renderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *Client) getTodoList(response http.ResponseWriter, request *http.Request) {
	todos := service.app.FindTodos()
	html := service.renderer.renderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *Client) addTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := service.app.AddTodo(task)
	if err != nil {
		html := service.renderer.renderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	html := service.renderer.renderTodoForm("", "")
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

func (client *Client) toggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = client.app.ToggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := client.app.FindTodo(id)
	html := client.renderer.renderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *Client) removeTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = service.app.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

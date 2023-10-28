package htmx

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (this *Client) getTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := this.app.FindTodos()
	html := this.renderer.renderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *Client) getTodoList(response http.ResponseWriter, request *http.Request) {
	todos := this.app.FindTodos()
	html := this.renderer.renderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *Client) addTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := this.app.AddTodo(task)
	if err != nil {
		html := this.renderer.renderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	html := this.renderer.renderTodoForm("", "")
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
	maybeId := url.Query().Get("id")
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

func (this *Client) toggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = this.app.ToggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := this.app.FindTodoById(id)
	html := this.renderer.renderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *Client) removeTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = this.app.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

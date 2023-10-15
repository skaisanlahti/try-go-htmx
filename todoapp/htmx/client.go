package htmx

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type Client struct {
	userAuthenticator *todoapp.UserAuthenticator
	todoWriter        *todoapp.TodoWriter
	htmlRenderer      *htmlRenderer
}

func NewClient(ua *todoapp.UserAuthenticator, tw *todoapp.TodoWriter) *Client {
	return &Client{ua, tw, newHtmlRenderer()}
}

func (client *Client) GetRegisterPage(response http.ResponseWriter, request *http.Request) {
	html := client.htmlRenderer.RenderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) RegisterUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := client.htmlRenderer.RenderRegisterForm(name, password, message)
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	}

	if name == "" {
		renderError("Username is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
		return
	}

	err := client.userAuthenticator.RegisterUser(name, password, response)
	if err != nil {
		renderError("Application error.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) GetLoginPage(response http.ResponseWriter, request *http.Request) {
	html := client.htmlRenderer.RenderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) GetLogoutPage(response http.ResponseWriter, request *http.Request) {
	html := client.htmlRenderer.RenderLogoutPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) LoginUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := client.htmlRenderer.RenderLoginForm(name, password, message)
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	}

	if name == "" {
		renderError("Username is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
		return
	}

	err := client.userAuthenticator.LoginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) LogoutUser(response http.ResponseWriter, request *http.Request) {
	err := client.userAuthenticator.LogoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) GetTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := client.todoWriter.TodoAccessor.FindTodos()
	html := client.htmlRenderer.RenderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) GetTodoList(response http.ResponseWriter, request *http.Request) {
	todos := client.todoWriter.TodoAccessor.FindTodos()
	html := client.htmlRenderer.RenderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) AddTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := client.todoWriter.AddTodo(task)
	if err != nil {
		html := client.htmlRenderer.RenderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	html := client.htmlRenderer.RenderTodoForm("", "")
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

func (client *Client) ToggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = client.todoWriter.ToggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := client.todoWriter.TodoAccessor.FindTodoById(id)
	html := client.htmlRenderer.RenderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) RemoveTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = client.todoWriter.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

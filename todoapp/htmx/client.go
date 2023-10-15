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
	UserAuthenticator *todoapp.UserAuthenticator
	TodoWriter        *todoapp.TodoWriter
	HtmlRenderer      *HtmlRenderer
}

func NewClient(ua *todoapp.UserAuthenticator, tw *todoapp.TodoWriter) *Client {
	return &Client{ua, tw, NewHtmlRenderer()}
}

func (client *Client) GetRegisterPage(response http.ResponseWriter, request *http.Request) {
	html := client.HtmlRenderer.RenderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) RegisterUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := client.HtmlRenderer.RenderRegisterForm(name, password, message)
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

	err := client.UserAuthenticator.RegisterUser(name, password, response)
	if err != nil {
		renderError("Application error.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) GetLoginPage(response http.ResponseWriter, request *http.Request) {
	html := client.HtmlRenderer.RenderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) GetLogoutPage(response http.ResponseWriter, request *http.Request) {
	html := client.HtmlRenderer.RenderLogoutPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) LoginUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := client.HtmlRenderer.RenderLoginForm(name, password, message)
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

	err := client.UserAuthenticator.LoginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) LogoutUser(response http.ResponseWriter, request *http.Request) {
	err := client.UserAuthenticator.LogoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) GetTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := client.TodoWriter.TodoAccessor.FindTodos()
	html := client.HtmlRenderer.RenderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) GetTodoList(response http.ResponseWriter, request *http.Request) {
	todos := client.TodoWriter.TodoAccessor.FindTodos()
	html := client.HtmlRenderer.RenderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) AddTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := client.TodoWriter.AddTodo(task)
	if err != nil {
		html := client.HtmlRenderer.RenderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
		return
	}

	html := client.HtmlRenderer.RenderTodoForm("", "")
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

	err = client.TodoWriter.ToggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := client.TodoWriter.TodoAccessor.FindTodoById(id)
	html := client.HtmlRenderer.RenderTodoItem(todo)
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

	err = client.TodoWriter.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

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
	AuthenticationService *todoapp.AuthenticationService
	TodoService           *todoapp.TodoService
	HtmlRenderer          *htmlRenderer
}

func CreateClient(as *todoapp.AuthenticationService, ts *todoapp.TodoService, router *http.ServeMux) {
	client := &Client{as, ts, createHtmlRenderer()}
	mapRoutes(router, client)
}

func (client *Client) getRegisterPage(response http.ResponseWriter, request *http.Request) {
	html := client.HtmlRenderer.RenderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) registerUser(response http.ResponseWriter, request *http.Request) {
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

	err := client.AuthenticationService.RegisterUser(name, password, response)
	if err != nil {
		renderError("Application error.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) getLoginPage(response http.ResponseWriter, request *http.Request) {
	html := client.HtmlRenderer.RenderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	html := client.HtmlRenderer.RenderLogoutPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) loginUser(response http.ResponseWriter, request *http.Request) {
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

	err := client.AuthenticationService.LoginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := client.AuthenticationService.LogoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}

func (client *Client) getTodoPage(response http.ResponseWriter, request *http.Request) {
	todos := client.TodoService.TodoStorage.FindTodos()
	html := client.HtmlRenderer.RenderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) getTodoList(response http.ResponseWriter, request *http.Request) {
	todos := client.TodoService.TodoStorage.FindTodos()
	html := client.HtmlRenderer.RenderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) addTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	err := client.TodoService.AddTodo(task)
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

func (client *Client) toggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = client.TodoService.ToggleTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, _ := client.TodoService.TodoStorage.FindTodoById(id)
	html := client.HtmlRenderer.RenderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) removeTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = client.TodoService.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

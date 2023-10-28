package htmx

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

//go:embed web/html/*.html
var templateFiles embed.FS

//go:embed web/dist/*
var assetFiles embed.FS

const (
	assetFilesRoot = "web/dist"
	assetPath      = "/dist/"
)

func useAssets(router *http.ServeMux) {
	assets, err := fs.Sub(assetFiles, assetFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	router.Handle(assetPath, http.StripPrefix(assetPath, http.FileServer(http.FS(assets))))
}

type renderer struct {
	loginPage    *template.Template
	logoutPage   *template.Template
	registerPage *template.Template
	todoPage     *template.Template
}

func newRenderer() *renderer {
	loginPage := template.Must(template.ParseFS(templateFiles, "web/html/login_page.html"))
	logoutPage := template.Must(template.ParseFS(templateFiles, "web/html/logout_page.html"))
	registerPage := template.Must(template.ParseFS(templateFiles, "web/html/register_page.html"))
	todoPage := template.Must(template.ParseFS(templateFiles, "web/html/todo_page.html"))
	return &renderer{loginPage, logoutPage, registerPage, todoPage}
}

type registerPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

func (renderer *renderer) renderRegisterPage() []byte {
	data := registerPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := renderer.registerPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *renderer) renderRegisterForm(name string, password string, errorMessage string) []byte {
	data := registerPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.registerPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type loginPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

func (renderer *renderer) renderLoginPage() []byte {
	data := loginPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *renderer) renderLoginForm(name string, password string, errorMessage string) []byte {
	data := loginPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type logoutPageData struct {
	LoggedIn bool
}

func (renderer *renderer) renderLogoutPage(loggedIn bool) []byte {
	data := logoutPageData{loggedIn}
	buffer := &bytes.Buffer{}
	err := renderer.logoutPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type todoPageData struct {
	Key   int64
	Task  string
	Todos []entity.Todo
	Error string
}

func (renderer *renderer) renderTodoPage(todos []entity.Todo) []byte {
	data := todoPageData{
		Key:   time.Now().UnixMilli(),
		Todos: todos,
	}

	buffer := &bytes.Buffer{}
	err := renderer.todoPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *renderer) renderTodoForm(task string, errorMessage string) []byte {
	data := todoPageData{
		Key:   time.Now().UnixMilli(),
		Task:  task,
		Error: errorMessage,
	}

	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *renderer) renderTodoList(todos []entity.Todo) []byte {
	data := todoPageData{Todos: todos}
	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "list", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *renderer) renderTodoItem(todo entity.Todo) []byte {
	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "item", todo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

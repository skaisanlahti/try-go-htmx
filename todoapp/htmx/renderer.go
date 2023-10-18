package htmx

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

//go:embed assets/dist/*
var assetFiles embed.FS

const (
	assetFilesRoot = "assets/dist"
	assetPath      = "/dist/"
)

func mapAssets(router *http.ServeMux) {
	assetFiles, err := fs.Sub(assetFiles, assetFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	router.Handle(assetPath, http.StripPrefix(assetPath, http.FileServer(http.FS(assetFiles))))
}

//go:embed assets/html/*.html
var templateFiles embed.FS

type htmlRenderer struct {
	loginPage    *template.Template
	logoutPage   *template.Template
	registerPage *template.Template
	todoPage     *template.Template
}

func newHtmlRenderer() *htmlRenderer {
	loginPage := template.Must(template.ParseFS(templateFiles, "assets/html/login_page.html"))
	logoutPage := template.Must(template.ParseFS(templateFiles, "assets/html/logout_page.html"))
	registerPage := template.Must(template.ParseFS(templateFiles, "assets/html/register_page.html"))
	todoPage := template.Must(template.ParseFS(templateFiles, "assets/html/todo_page.html"))
	return &htmlRenderer{loginPage, logoutPage, registerPage, todoPage}
}

type registerPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

func (renderer *htmlRenderer) RenderRegisterPage() []byte {
	data := registerPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := renderer.registerPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmlRenderer) RenderRegisterForm(name string, password string, errorMessage string) []byte {
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

func (renderer *htmlRenderer) RenderLoginPage() []byte {
	data := loginPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.Execute(buffer, data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmlRenderer) RenderLoginForm(name string, password string, errorMessage string) []byte {
	data := loginPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := renderer.loginPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmlRenderer) RenderLogoutPage() []byte {
	buffer := &bytes.Buffer{}
	err := renderer.logoutPage.Execute(buffer, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type TodoPageData struct {
	Key   int64
	Task  string
	Todos []todoapp.Todo
	Error string
}

func (renderer *htmlRenderer) RenderTodoPage(todos []todoapp.Todo) []byte {
	data := TodoPageData{
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

func (renderer *htmlRenderer) RenderTodoForm(task string, errorMessage string) []byte {
	data := TodoPageData{
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

func (renderer *htmlRenderer) RenderTodoList(todos []todoapp.Todo) []byte {
	data := TodoPageData{Todos: todos}
	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "list", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmlRenderer) RenderTodoItem(todo todoapp.Todo) []byte {
	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "item", todo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

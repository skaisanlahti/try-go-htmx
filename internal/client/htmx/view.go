package htmx

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

//go:embed web/html/*.html
var templateFiles embed.FS

type view struct {
	loginPage    *template.Template
	logoutPage   *template.Template
	registerPage *template.Template
	todoPage     *template.Template
}

func newView() *view {
	loginPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/login_page.html"))
	logoutPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/logout_page.html"))
	registerPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/register_page.html"))
	todoPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/todo_page.html"))
	return &view{loginPage, logoutPage, registerPage, todoPage}
}

type registerPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

func (this *view) renderRegisterPage() []byte {
	data := registerPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := this.registerPage.ExecuteTemplate(buffer, "page", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (this *view) renderRegisterForm(name string, password string, errorMessage string) []byte {
	data := registerPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := this.registerPage.ExecuteTemplate(buffer, "form", data)
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

func (this *view) renderLoginPage() []byte {
	data := loginPageData{Key: time.Now().UnixMilli()}
	buffer := &bytes.Buffer{}
	err := this.loginPage.ExecuteTemplate(buffer, "page", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (this *view) renderLoginForm(name string, password string, errorMessage string) []byte {
	data := loginPageData{Key: time.Now().UnixMilli(), Name: name, Password: password, Error: errorMessage}
	buffer := &bytes.Buffer{}
	err := this.loginPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

type logoutPageData struct {
	LoggedIn bool
}

func (this *view) renderLogoutPage(loggedIn bool) []byte {
	data := logoutPageData{loggedIn}
	buffer := &bytes.Buffer{}
	err := this.logoutPage.ExecuteTemplate(buffer, "page", data)
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

func (this *view) renderTodoPage(todos []entity.Todo) []byte {
	data := todoPageData{
		Key:   time.Now().UnixMilli(),
		Todos: todos,
	}

	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "page", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (this *view) renderTodoForm(task string, errorMessage string) []byte {
	data := todoPageData{
		Key:   time.Now().UnixMilli(),
		Task:  task,
		Error: errorMessage,
	}

	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "form", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (this *view) renderTodoList(todos []entity.Todo) []byte {
	data := todoPageData{Todos: todos}
	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "list", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (this *view) renderTodoItem(todo entity.Todo) []byte {
	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "item", todo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

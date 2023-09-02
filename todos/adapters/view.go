package adapters

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/models"
)

//go:embed templates/todo_page.html
var templateFiles embed.FS

func newTodoPageTemplate() *template.Template {
	return template.Must(template.ParseFS(templateFiles, "templates/todo_page.html"))
}

type View struct {
	todoPageTemplate *template.Template
}

func NewView() *View {
	return &View{newTodoPageTemplate()}
}

func (this *View) RenderPage(response http.ResponseWriter, data models.TodoPage) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.todoPageTemplate.Execute(response, data)
}

func (this *View) RenderForm(response http.ResponseWriter, data models.TodoPage) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.todoPageTemplate.ExecuteTemplate(response, "form", data)
}

func (this *View) RenderList(response http.ResponseWriter, data models.TodoPage) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.todoPageTemplate.ExecuteTemplate(response, "list", data)
}

func (this *View) RenderItem(response http.ResponseWriter, data models.Todo) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.todoPageTemplate.ExecuteTemplate(response, "item", data)
}

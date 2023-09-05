package procedures

import (
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/data"
)

type HtmlRenderer struct {
	todoPageTemplate *template.Template
}

func NewHtmlRenderer(todoPageTemplate *template.Template) *HtmlRenderer {
	return &HtmlRenderer{todoPageTemplate}
}

func (this *HtmlRenderer) RenderPage(response http.ResponseWriter, todos []data.Todo) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	templateData := data.NewTodoPage()
	templateData.Todos = todos
	return this.todoPageTemplate.Execute(response, templateData)
}

func (this *HtmlRenderer) RenderList(response http.ResponseWriter, todos []data.Todo) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	templateData := data.NewTodoPage()
	templateData.Todos = todos
	return this.todoPageTemplate.ExecuteTemplate(response, "list", templateData)
}

func (this *HtmlRenderer) RenderEmptyForm(response http.ResponseWriter) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	templateData := data.NewTodoPage()
	return this.todoPageTemplate.ExecuteTemplate(response, "form", templateData)
}

func (this *HtmlRenderer) RenderErrorForm(response http.ResponseWriter, errorMessage string) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	templateData := data.NewTodoPage()
	templateData.Error = errorMessage
	return this.todoPageTemplate.ExecuteTemplate(response, "form", templateData)
}

func (this *HtmlRenderer) RenderItem(response http.ResponseWriter, todo data.Todo) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.todoPageTemplate.ExecuteTemplate(response, "item", todo)
}

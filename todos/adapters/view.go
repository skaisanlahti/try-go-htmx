package adapters

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/models"
)

//go:embed page.html
var templateFiles embed.FS

func newTaskPageTemplate() *template.Template {
	return template.Must(template.ParseFS(templateFiles, "page.html"))
}

type View struct {
	pageTemplate *template.Template
}

func NewView() *View {
	return &View{newTaskPageTemplate()}
}

func (this *View) RenderPage(response http.ResponseWriter, data models.TodoPage) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.pageTemplate.Execute(response, data)
}

func (this *View) RenderForm(response http.ResponseWriter, data models.TodoPage) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.pageTemplate.ExecuteTemplate(response, "form", data)

}

func (this *View) RenderList(response http.ResponseWriter, data models.TodoPage) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.pageTemplate.ExecuteTemplate(response, "list", data)
}

func (this *View) RenderItem(response http.ResponseWriter, data models.Todo) error {
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	return this.pageTemplate.ExecuteTemplate(response, "item", data)
}

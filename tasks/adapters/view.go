package adapters

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/tasks/models"
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

func (this *View) RenderPage(response http.ResponseWriter, data models.TaskPage) error {
	return this.pageTemplate.Execute(response, data)
}

func (this *View) RenderMain(response http.ResponseWriter, data models.TaskPage) error {
	return this.pageTemplate.ExecuteTemplate(response, "main", data)

}

func (this *View) RenderList(response http.ResponseWriter, data models.TaskPage) error {
	return this.pageTemplate.ExecuteTemplate(response, "list", data)
}

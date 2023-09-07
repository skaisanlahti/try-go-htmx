package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var templateFiles embed.FS

type HtmlTemplates struct {
	TodoPage *template.Template
}

func ParseTemplates() *HtmlTemplates {
	todoPage := template.Must(template.ParseFS(templateFiles, "todo_page.html"))
	return &HtmlTemplates{todoPage}
}

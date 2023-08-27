package todos

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
)

const (
	templateMain string = "main"
	templateList string = "list"
)

//go:embed page.html
var templateFiles embed.FS

type todoRenderer interface {
	renderPage(response http.ResponseWriter, data todoPageData)
	renderMain(response http.ResponseWriter, data todoPageData)
	renderList(response http.ResponseWriter, data todoPageData)
}

type todoView struct {
	todoPage *template.Template
}

func newTodoView() *todoView {
	todoPage := template.Must(template.ParseFS(templateFiles, "page.html"))
	return &todoView{
		todoPage: todoPage,
	}
}

func (this *todoView) renderPage(response http.ResponseWriter, data todoPageData) {
	buffer := &bytes.Buffer{}
	err := this.todoPage.Execute(buffer, data)
	if err != nil {
		log.Print(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(buffer.Bytes())
}

func (this *todoView) renderMain(response http.ResponseWriter, data todoPageData) {
	buffer := &bytes.Buffer{}
	if err := this.todoPage.ExecuteTemplate(buffer, templateMain, data); err != nil {
		log.Print(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(buffer.Bytes())
}

func (this *todoView) renderList(response http.ResponseWriter, data todoPageData) {
	buffer := &bytes.Buffer{}
	if err := this.todoPage.ExecuteTemplate(buffer, templateList, data); err != nil {
		log.Print(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(buffer.Bytes())
}

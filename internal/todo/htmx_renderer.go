package todo

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"time"
)

type htmxRenderer struct {
	todoPage *template.Template
}

func NewHtmxRenderer(files embed.FS) *htmxRenderer {
	todoPage := template.Must(template.ParseFS(files, "html/todo_page.html"))
	return &htmxRenderer{todoPage}
}

type registerPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

type todoPageData struct {
	Key   int64
	Task  string
	Todos []todo
	Error string
}

func (renderer *htmxRenderer) renderTodoPage(todos []todo) []byte {
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

func (renderer *htmxRenderer) renderTodoForm(task string, errorMessage string) []byte {
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

func (renderer *htmxRenderer) renderTodoList(todos []todo) []byte {
	data := todoPageData{Todos: todos}
	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "list", data)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

func (renderer *htmxRenderer) renderTodoItem(todo todo) []byte {
	buffer := &bytes.Buffer{}
	err := renderer.todoPage.ExecuteTemplate(buffer, "item", todo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

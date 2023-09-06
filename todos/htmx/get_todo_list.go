package htmx

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/domain"
)

type GetTodoListRepository interface {
	GetTodos() []domain.Todo
}

type GetTodoListView interface {
	RenderTodoList(todos []domain.Todo) []byte
}

type GetTodoListHandler struct {
	repository GetTodoListRepository
	view       GetTodoListView
}

func NewGetTodoListHandler(repository GetTodoListRepository, view GetTodoListView) *GetTodoListHandler {
	return &GetTodoListHandler{repository, view}
}

func (this *GetTodoListHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	todos := this.repository.GetTodos()
	html := this.view.RenderTodoList(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.Write(html)
}

type TodoList struct {
	Todos []domain.Todo
}

type HtmxGetTodoListView struct {
	todoPage *template.Template
}

func NewHtmxGetTodoListView(todoPageTemplate *template.Template) *HtmxGetTodoListView {
	return &HtmxGetTodoListView{todoPageTemplate}
}

func (this *HtmxGetTodoListView) RenderTodoList(todos []domain.Todo) []byte {
	templateData := TodoList{Todos: todos}
	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "list", templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

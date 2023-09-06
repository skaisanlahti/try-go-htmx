package htmx

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/todos/domain"
)

type GetTodoPageRepository interface {
	GetTodos() []domain.Todo
}

type GetTodoPageView interface {
	RenderTodoPage(todos []domain.Todo) []byte
}

type GetTodoPageHandler struct {
	repository GetTodoPageRepository
	view       GetTodoPageView
}

func NewGetTodoPageHandler(repository GetTodoPageRepository, view GetTodoPageView) *GetTodoPageHandler {
	return &GetTodoPageHandler{repository, view}
}

func (handler *GetTodoPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	todos := handler.repository.GetTodos()
	html := handler.view.RenderTodoPage(todos)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

type TodoPage struct {
	Key   int64
	Task  string
	Error string
	Todos []domain.Todo
}

type HtmxGetTodoPageView struct {
	todoPage *template.Template
}

func NewHtmxGetTodoPageView(todoPage *template.Template) *HtmxGetTodoPageView {
	return &HtmxGetTodoPageView{todoPage}
}

func (view *HtmxGetTodoPageView) RenderTodoPage(todos []domain.Todo) []byte {
	templateData := TodoPage{Key: time.Now().UnixMilli(), Todos: todos}
	buffer := &bytes.Buffer{}
	err := view.todoPage.Execute(buffer, templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

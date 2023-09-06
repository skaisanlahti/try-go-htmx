package htmx

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/todos/domain"
)

type AddTodoRepository interface {
	AddTodo(todo domain.Todo) error
}

type AddTodoView interface {
	RenderTodoForm(task string, taskError string) []byte
}

type AddTodoHandler struct {
	repository AddTodoRepository
	view       AddTodoView
}

func NewAddTodoHandler(repository AddTodoRepository, view AddTodoView) *AddTodoHandler {
	return &AddTodoHandler{repository, view}
}

func (this *AddTodoHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	newTodo, err := domain.NewTodo(task)
	if err != nil {
		html := this.view.RenderTodoForm(task, err.Error())
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.Write(html)
		return
	}

	if err := this.repository.AddTodo(newTodo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	html := this.view.RenderTodoForm("", "")
	response.Header().Add("HX-Trigger", "GetTodoList")
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.Write(html)
}

type TodoForm struct {
	Key   int64
	Task  string
	Error string
}

type HtmxAddTodoView struct {
	todoPage *template.Template
}

func NewHtmxAddTodoView(view *template.Template) *HtmxAddTodoView {
	return &HtmxAddTodoView{view}
}

func (this *HtmxAddTodoView) RenderTodoForm(task string, taskError string) []byte {
	templateData := TodoForm{Key: time.Now().UnixMilli(), Task: task, Error: taskError}
	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "form", templateData)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

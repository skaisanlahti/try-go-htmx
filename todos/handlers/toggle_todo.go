package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/domain"
)

type ToggleTodoRepository interface {
	GetTodoById(id int) (domain.Todo, error)
	UpdateTodo(todo domain.Todo) error
}

type ToggleTodoView interface {
	RenderTodoItem(todo domain.Todo) []byte
}

type ToggleTodoHandler struct {
	repository ToggleTodoRepository
	view       ToggleTodoView
}

func NewToggleTodoHandler(
	repository ToggleTodoRepository,
	view ToggleTodoView,
) *ToggleTodoHandler {
	return &ToggleTodoHandler{repository, view}
}

func (this *ToggleTodoHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := this.repository.GetTodoById(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	newTodo := domain.ToggleTodo(todo)

	if err = this.repository.UpdateTodo(newTodo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	html := this.view.RenderTodoItem(todo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.Write(html)
}

type TemplateToggleTodoView struct {
	todoPage *template.Template
}

func NewTemplateToggleTodoView(view *template.Template) *TemplateToggleTodoView {
	return &TemplateToggleTodoView{view}
}

func (this *TemplateToggleTodoView) RenderTodoItem(todo domain.Todo) []byte {
	buffer := &bytes.Buffer{}
	err := this.todoPage.ExecuteTemplate(buffer, "item", todo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

package htmx

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

func (handler *ToggleTodoHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := handler.repository.GetTodoById(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	newTodo := domain.ToggleTodo(todo)

	if err = handler.repository.UpdateTodo(newTodo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	html := handler.view.RenderTodoItem(newTodo)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

type HtmxToggleTodoView struct {
	todoPage *template.Template
}

func NewHtmxToggleTodoView(todoPage *template.Template) *HtmxToggleTodoView {
	return &HtmxToggleTodoView{todoPage}
}

func (view *HtmxToggleTodoView) RenderTodoItem(todo domain.Todo) []byte {
	buffer := &bytes.Buffer{}
	err := view.todoPage.ExecuteTemplate(buffer, "item", todo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return buffer.Bytes()
}

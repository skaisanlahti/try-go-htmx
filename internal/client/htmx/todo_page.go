package htmx

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type todoPageData struct {
	Key          int64
	TodoListId   int
	TodoListName string
	Task         string
	Todos        []entity.Todo
	Error        string
}

type todoPageController struct {
	todo *todo.TodoService
	*defaultRenderer
}

func newTodoPageController(todo *todo.TodoService) *todoPageController {
	todoPageTemplate := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/todo_page.html"))
	return &todoPageController{todo, newDefaultRenderer(todoPageTemplate)}
}

var (
	ErrListIdMissing   = errors.New("List id not found in query.")
	ErrListIdNotNumber = errors.New("List id not a number.")
)

func extractListId(url *url.URL) (int, error) {
	maybeId := url.Query().Get("listid")
	if maybeId == "" {
		return 0, ErrListIdMissing
	}

	id, err := strconv.Atoi(maybeId)
	if err != nil {
		log.Println(err.Error())
		return 0, ErrListIdNotNumber
	}

	return id, nil
}

func (this *todoPageController) page(response http.ResponseWriter, request *http.Request) {
	listId, err := extractListId(request.URL)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(response, request, "/htmx/todo-lists", http.StatusSeeOther)
		return
	}

	list, err := this.todo.FindListById(listId)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(response, request, "/htmx/todo-lists", http.StatusSeeOther)
		return
	}

	this.render(response, "page", todoPageData{
		Key:          newRenderKey(),
		TodoListId:   listId,
		TodoListName: list.Name,
		Todos:        this.todo.FindTodosByListId(listId),
	}, nil)

}

func (this *todoPageController) todos(response http.ResponseWriter, request *http.Request) {
	listId, err := extractListId(request.URL)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	this.render(response, "list", todoPageData{
		TodoListId: listId,
		Todos:      this.todo.FindTodosByListId(listId),
	}, nil)

}

func (this *todoPageController) addTodo(response http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task")
	maybeListId := request.FormValue("listId")
	listId, err := strconv.Atoi(maybeListId)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = this.todo.AddTodo(task, listId); err != nil {
		this.render(response, "form", todoPageData{
			Key:   newRenderKey(),
			TodoListId: listId,
			Task:  task,
			Error: err.Error(),
		}, nil)
		return
	}

	this.render(response, "form", todoPageData{
		Key: newRenderKey(),
		TodoListId: listId,
	}, extraHeaders{
		"HX-Trigger": "GetTodos",
	})
}

var (
	ErrTodoIdMissing   = errors.New("Todo id not found in query.")
	ErrTodoIdNotNumber = errors.New("Todo id not a number.")
)

func extractTodoId(url *url.URL) (int, error) {
	maybeId := url.Query().Get("id")
	if maybeId == "" {
		return 0, ErrTodoIdMissing
	}

	id, err := strconv.Atoi(maybeId)
	if err != nil {
		log.Println(err.Error())
		return 0, ErrTodoIdNotNumber
	}

	return id, nil
}

func (this *todoPageController) toggleTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := this.todo.ToggleTodo(id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	this.render(response, "item", todo, nil)
}

func (this *todoPageController) removeTodo(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = this.todo.RemoveTodo(id); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

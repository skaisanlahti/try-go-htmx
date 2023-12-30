package htmx

import (
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type todoListPageData struct {
	Key       int64
	Name      string
	TodoLists []entity.TodoList
	Error     string
}

type todoListPageController struct {
	todo *todo.TodoService
	*defaultRenderer
}

func newTodoListPageController(todo *todo.TodoService) *todoListPageController {
	todoListPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/todo_list_page.html"))
	return &todoListPageController{todo, newDefaultRenderer(todoListPage)}
}

func (this *todoListPageController) page(response http.ResponseWriter, request *http.Request) {
	user, ok := extractUserFromContext(request)
	if !ok {
		http.Error(response, "User not found.", http.StatusBadRequest)
		return
	}

	this.render(response, "page", todoListPageData{
		Key:       newRenderKey(),
		TodoLists: this.todo.FindListsByUserId(user.Id),
	}, nil)
}

func (this *todoListPageController) lists(response http.ResponseWriter, request *http.Request) {
	user, ok := extractUserFromContext(request)
	if !ok {
		http.Error(response, "User not found.", http.StatusBadRequest)
		return
	}

	this.render(response, "list", todoListPageData{
		Key:       newRenderKey(),
		TodoLists: this.todo.FindListsByUserId(user.Id),
	}, nil)
}

func (this *todoListPageController) addList(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	user, ok := extractUserFromContext(request)
	if !ok {
		http.Error(response, "User not found.", http.StatusBadRequest)
		return
	}

	if _, err := this.todo.AddList(name, user.Id); err != nil {
		this.render(response, "form", todoListPageData{
			Key:   newRenderKey(),
			Name:  name,
			Error: err.Error(),
		}, nil)
		return
	}

	this.render(response, "form", todoListPageData{
		Key: newRenderKey(),
	}, extraHeaders{
		"HX-Trigger": "GetLists",
	})
}

func (this *todoListPageController) removeList(response http.ResponseWriter, request *http.Request) {
	listId, err := extractListId(request.URL)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = this.todo.RemoveList(listId); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

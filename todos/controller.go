package todos

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	c "github.com/skaisanlahti/test-go/common"
)

type todoController struct {
	repository todoRepository
	view       todoRenderer
}

func newTodoController(database *sql.DB) *todoController {
	return &todoController{
		repository: newSqlTodoRepository(database),
		view:       newTodoView(),
	}
}

func (this *todoController) todoPage(response http.ResponseWriter, request *http.Request) error {
	data, err := newTodoPageData(this.repository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.view.renderPage(response, data)
	return nil
}

func (this *todoController) addTodo(response http.ResponseWriter, request *http.Request) error {
	task := request.FormValue("task")
	if task == "" {
		data, err := newTodoPageData(this.repository)
		if err != nil {
			return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
		}

		data.Error = "Task can't be empty."
		this.view.renderMain(response, data)
		return nil
	}

	newTodo := todoRecord{Task: task, Done: false}
	if err := this.repository.addTodo(newTodo); err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	data, err := newTodoPageData(this.repository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.view.renderMain(response, data)
	return nil
}

func (this *todoController) toggleTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusBadRequest).Log(err)
	}

	todo, err := this.repository.getTodoByID(id)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	todo.Done = !todo.Done
	err = this.repository.updateTodo(todo)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	data, err := newTodoPageData(this.repository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.view.renderList(response, data)
	return nil
}

func (this *todoController) removeTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusBadRequest).Log(err)
	}

	err = this.repository.removeTodo(id)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	data, err := newTodoPageData(this.repository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.view.renderList(response, data)
	return nil
}

func extractTodoID(url *url.URL) (int, error) {
	values := url.Query()
	idStr := values.Get("id")
	if idStr == "" {
		return 0, errors.New("Todo ID not found in query")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

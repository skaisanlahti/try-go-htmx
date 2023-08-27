package todos

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	c "github.com/skaisanlahti/try-go-htmx/common"
)

type todoController struct {
	todoRepository repository[todoRecord]
	todoView       todoRenderer
}

func newTodoController(database *sql.DB) *todoController {
	return &todoController{
		todoRepository: newSqlTodoRepository(database),
		todoView:       newTodoView(),
	}
}

func (this *todoController) todoPage(response http.ResponseWriter, request *http.Request) error {
	data, err := newTodoPageData(this.todoRepository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.todoView.renderPage(response, data)
	return nil
}

func (this *todoController) addTodo(response http.ResponseWriter, request *http.Request) error {
	task := request.FormValue("task")
	if task == "" {
		data, err := newTodoPageData(this.todoRepository)
		if err != nil {
			return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
		}

		data.Error = "Task can't be empty."
		this.todoView.renderMain(response, data)
		return nil
	}

	newTodo := todoRecord{Task: task, Done: false}
	if err := this.todoRepository.add(newTodo); err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	data, err := newTodoPageData(this.todoRepository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.todoView.renderMain(response, data)
	return nil
}

func (this *todoController) toggleTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusBadRequest).Log(err)
	}

	todo, err := this.todoRepository.find(id)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusNotFound).Log(err)
	}

	todo.Done = !todo.Done
	err = this.todoRepository.update(todo)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	data, err := newTodoPageData(this.todoRepository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.todoView.renderList(response, data)
	return nil
}

func (this *todoController) removeTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusBadRequest).Log(err)
	}

	err = this.todoRepository.remove(id)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	data, err := newTodoPageData(this.todoRepository)
	if err != nil {
		return c.NewRouteError(err.Error(), http.StatusInternalServerError).Log(err)
	}

	this.todoView.renderList(response, data)
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

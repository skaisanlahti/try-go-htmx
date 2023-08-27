package todos

import (
	"time"
)

type todoPageData struct {
	Todos []todoRecord
	Task  string
	Error string
	Key   int64 // use as cache key or name attribute to bypass caching in browser
}

func newTodoPageData(todoClient todoRepository) (todoPageData, error) {
	todos, err := todoClient.getAllTodos()
	if err != nil {
		return todoPageData{}, err
	}

	data := todoPageData{Key: time.Now().UnixMilli(), Todos: todos}
	return data, nil
}

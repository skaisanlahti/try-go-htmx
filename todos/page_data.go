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

func newTodoPageData(todoRepository repository[todoRecord]) (todoPageData, error) {
	todos, err := todoRepository.list()
	if err != nil {
		return todoPageData{}, err
	}

	data := todoPageData{Key: time.Now().UnixMilli(), Todos: todos}
	return data, nil
}

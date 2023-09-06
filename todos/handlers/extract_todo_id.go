package handlers

import (
	"errors"
	"net/url"
	"strconv"
)

func extractTodoId(url *url.URL) (int, error) {
	values := url.Query()
	idStr := values.Get("id")
	if idStr == "" {
		return 0, errors.New("Task Id not found in query")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

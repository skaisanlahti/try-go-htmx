package handlers

import (
	"errors"
	"log"
	"net/url"
	"strconv"
)

func extractTodoId(url *url.URL) (int, error) {
	values := url.Query()
	maybeId := values.Get("id")
	if maybeId == "" {
		return 0, errors.New("Todo id not found in query.")
	}

	id, err := strconv.Atoi(maybeId)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return id, nil
}

package core

import (
	"net/http"
)

func (this *Service) RemoveTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = this.Database.RemoveTodo(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.WriteHeader(http.StatusOK)
	return nil
}

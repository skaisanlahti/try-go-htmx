package core

import (
	"net/http"
)

func (this *Service) RemoveTask(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTaskID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = this.Database.RemoveTask(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.WriteHeader(http.StatusOK)
	return nil
}

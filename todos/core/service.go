package core

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/todos/adapters"
	"github.com/skaisanlahti/try-go-htmx/todos/ports"
)

type Service struct {
	Database ports.Database
	View     ports.View
}

func NewService(database *sql.DB) *Service {
	return &Service{adapters.NewDatabase(database), adapters.NewView()}
}

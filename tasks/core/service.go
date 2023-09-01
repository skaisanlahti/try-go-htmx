package core

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/tasks/adapters"
	"github.com/skaisanlahti/try-go-htmx/tasks/ports"
)

type Service struct {
	Database ports.Database
	View     ports.View
}

func NewService(database *sql.DB) *Service {
	return &Service{adapters.NewDatabase(database), adapters.NewView()}
}

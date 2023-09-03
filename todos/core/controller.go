package core

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/todos/adapters"
	"github.com/skaisanlahti/try-go-htmx/todos/ports"
)

type Controller struct {
	Database ports.Database
	View     ports.View
}

func NewController(database *sql.DB) *Controller {
	return &Controller{adapters.NewDatabase(database), adapters.NewView()}
}

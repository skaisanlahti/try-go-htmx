package ports

import "github.com/skaisanlahti/try-go-htmx/tasks/models"

type Database interface {
	GetTasks() ([]models.Task, error)
	GetTaskByID(id int) (models.Task, error)
	AddTask(task models.Task) error
	UpdateTask(task models.Task) error
	RemoveTask(id int) error
}

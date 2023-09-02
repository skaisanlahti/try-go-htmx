package adapters

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/tasks/models"
)

type Database struct {
	query *Query
}

func NewDatabase(database *sql.DB) *Database {
	return &Database{NewQuery(database)}
}

func (this *Database) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	query := this.query.Prepare(`SELECT * FROM todos ORDER BY task ASC`)
	rows, err := query.Query()
	if err != nil {
		return tasks, err
	}

	defer rows.Close()
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Task, &task.Done); err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (this *Database) GetTaskByID(id int) (models.Task, error) {
	var task models.Task
	query := this.query.Prepare(`SELECT * FROM todos WHERE todo_id = $1`)
	if err := query.QueryRow(id).Scan(&task.Id, &task.Task, &task.Done); err != nil {
		return task, err
	}

	return task, nil
}

func (this *Database) AddTask(task models.Task) error {
	insert := this.query.Prepare(`INSERT INTO todos (task, done) VALUES ($1, $2) RETURNING todo_id`)
	if _, err := insert.Exec(&task.Task, &task.Done); err != nil {
		return err
	}

	return nil
}

func (this *Database) UpdateTask(task models.Task) error {
	update := this.query.Prepare(`UPDATE todos SET task = $2, done = $3 WHERE todo_id = $1`)
	if _, err := update.Exec(&task.Id, &task.Task, &task.Done); err != nil {
		return err
	}

	return nil
}

func (this *Database) RemoveTask(id int) error {
	delete := this.query.Prepare(`DELETE FROM "Todos" WHERE "Id" = $1`)
	if _, err := delete.Exec(id); err != nil {
		return err
	}

	return nil

}

package adapters

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/todos/models"
	"github.com/skaisanlahti/try-go-htmx/todos/ports"
)

type Database struct {
	query ports.Query
}

func NewDatabase(database *sql.DB) *Database {
	return &Database{NewQuery(database)}
}

func (this *Database) GetTodos() ([]models.Todo, error) {
	var todos []models.Todo
	query := this.query.Prepare(`SELECT * FROM "Todos" ORDER BY "Task" ASC`)
	rows, err := query.Query()
	if err != nil {
		return todos, err
	}

	defer rows.Close()
	for rows.Next() {
		var task models.Todo
		if err := rows.Scan(&task.Id, &task.Task, &task.Done); err != nil {
			return todos, err
		}

		todos = append(todos, task)
	}

	return todos, nil
}

func (this *Database) GetTodoByID(id int) (models.Todo, error) {
	var todo models.Todo
	query := this.query.Prepare(`SELECT * FROM "Todos" WHERE "Id" = $1`)
	if err := query.QueryRow(id).Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return todo, err
	}

	return todo, nil
}

func (this *Database) AddTodo(todo models.Todo) error {
	insert := this.query.Prepare(`INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2) RETURNING "Id"`)
	if _, err := insert.Exec(&todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

func (this *Database) UpdateTodo(todo models.Todo) error {
	update := this.query.Prepare(`UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`)
	if _, err := update.Exec(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

func (this *Database) RemoveTodo(id int) error {
	delete := this.query.Prepare(`DELETE FROM "Todos" WHERE "Id" = $1`)
	if _, err := delete.Exec(id); err != nil {
		return err
	}

	return nil

}

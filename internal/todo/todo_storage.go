package todo

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type todoStorage struct {
	database *sql.DB
}

func newTodoStorage(database *sql.DB) *todoStorage {
	return &todoStorage{database}
}

func (this *todoStorage) findTodos() []entity.Todo {
	var todos []entity.Todo
	query := `SELECT * FROM "Todos" ORDER BY "Task" ASC`
	rows, err := this.database.Query(query)
	if err != nil {
		log.Println(err.Error())
		return todos
	}

	defer rows.Close()
	for rows.Next() {
		var todo entity.Todo
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
			log.Println(err.Error())
			return todos
		}

		todos = append(todos, todo)
	}

	return todos
}

func (this *todoStorage) findTodoById(id int) (entity.Todo, error) {
	var todo entity.Todo
	query := `SELECT * FROM "Todos" WHERE "Id" = $1`
	row := this.database.QueryRow(query, id)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func (this *todoStorage) insertTodo(todo entity.Todo) error {
	query := `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2)`
	if _, err := this.database.Exec(query, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (this *todoStorage) updateTodo(todo entity.Todo) error {
	query := `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`
	if _, err := this.database.Exec(query, &todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (this *todoStorage) deleteTodo(id int) error {
	query := `DELETE FROM "Todos" WHERE "Id" = $1`
	if _, err := this.database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

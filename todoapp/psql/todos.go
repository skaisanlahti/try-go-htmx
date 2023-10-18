package psql

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type TodoStorage struct {
	Database *sql.DB
}

func CreateTodoStorage(db *sql.DB) *TodoStorage {
	return &TodoStorage{db}
}

func (storage *TodoStorage) FindTodos() []todoapp.Todo {
	var todos []todoapp.Todo
	sql := `SELECT * FROM "Todos" ORDER BY "Task" ASC`
	rows, err := storage.Database.Query(sql)
	if err != nil {
		log.Println(err.Error())
		return todos
	}

	defer rows.Close()
	for rows.Next() {
		var todo todoapp.Todo
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
			log.Println(err.Error())
			return todos
		}

		todos = append(todos, todo)
	}

	return todos
}

func (storage *TodoStorage) FindTodoById(id int) (todoapp.Todo, error) {
	var todo todoapp.Todo
	sql := `SELECT * FROM "Todos" WHERE "Id" = $1`
	row := storage.Database.QueryRow(sql, id)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func (storage *TodoStorage) AddTodo(todo todoapp.Todo) error {
	sql := `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2)`
	if _, err := storage.Database.Exec(sql, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *TodoStorage) UpdateTodo(todo todoapp.Todo) error {
	sql := `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`
	if _, err := storage.Database.Exec(sql, &todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *TodoStorage) RemoveTodo(id int) error {
	sql := `DELETE FROM "Todos" WHERE "Id" = $1`
	if _, err := storage.Database.Exec(sql, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

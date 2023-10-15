package psql

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type TodoAccessor struct {
	Database *sql.DB
}

func NewTodoAccessor(db *sql.DB) *TodoAccessor {
	return &TodoAccessor{db}
}

func (accessor *TodoAccessor) FindTodos() []todoapp.Todo {
	var todos []todoapp.Todo
	sql := `SELECT * FROM "Todos" ORDER BY "Task" ASC`
	rows, err := accessor.Database.Query(sql)
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

func (accessor *TodoAccessor) FindTodoById(id int) (todoapp.Todo, error) {
	var todo todoapp.Todo
	sql := `SELECT * FROM "Todos" WHERE "Id" = $1`
	row := accessor.Database.QueryRow(sql, id)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func (accessor *TodoAccessor) AddTodo(todo todoapp.Todo) error {
	sql := `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2)`
	if _, err := accessor.Database.Exec(sql, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (accessor *TodoAccessor) UpdateTodo(todo todoapp.Todo) error {
	sql := `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`
	if _, err := accessor.Database.Exec(sql, &todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (accessor *TodoAccessor) RemoveTodo(id int) error {
	sql := `DELETE FROM "Todos" WHERE "Id" = $1`
	if _, err := accessor.Database.Exec(sql, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

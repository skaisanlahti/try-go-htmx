package services

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/todos/data"
)

type Preparer interface {
	Prepare(query string) *sql.Stmt
}

type PostgreSqlStorage struct {
	preparer Preparer
}

func NewPostgreSqlStorage(preparer Preparer) *PostgreSqlStorage {
	return &PostgreSqlStorage{preparer}
}

func (this *PostgreSqlStorage) FindTodos() ([]data.Todo, error) {
	var todos []data.Todo
	finder := this.preparer.Prepare(`SELECT * FROM "Todos" ORDER BY "Task" ASC`)
	rows, err := finder.Query()
	if err != nil {
		return todos, err
	}

	defer rows.Close()
	for rows.Next() {
		var todo data.Todo
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (this *PostgreSqlStorage) FindTodoByID(id int) (data.Todo, error) {
	var todo data.Todo
	finder := this.preparer.Prepare(`SELECT * FROM "Todos" WHERE "Id" = $1`)
	if err := finder.QueryRow(id).Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return todo, err
	}

	return todo, nil
}

func (this *PostgreSqlStorage) AddTodo(todo data.Todo) error {
	adder := this.preparer.Prepare(`INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2) RETURNING "Id"`)
	if _, err := adder.Exec(&todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

func (this *PostgreSqlStorage) UpdateTodo(todo data.Todo) error {
	updater := this.preparer.Prepare(`UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`)
	if _, err := updater.Exec(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

func (this *PostgreSqlStorage) RemoveTodo(id int) error {
	remover := this.preparer.Prepare(`DELETE FROM "Todos" WHERE "Id" = $1`)
	if _, err := remover.Exec(id); err != nil {
		return err
	}

	return nil

}

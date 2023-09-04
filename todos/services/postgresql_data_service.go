package services

import (
	"github.com/skaisanlahti/try-go-htmx/todos/data"
	"github.com/skaisanlahti/try-go-htmx/todos/interfaces"
)

type PostgreSqlDataService struct {
	queryService interfaces.QueryService
}

func NewPostgreSqlDataService(queryService interfaces.QueryService) *PostgreSqlDataService {
	return &PostgreSqlDataService{queryService}
}

func (this *PostgreSqlDataService) FindTodos() ([]data.Todo, error) {
	var todos []data.Todo
	query := this.queryService.Prepare(`SELECT * FROM "Todos" ORDER BY "Task" ASC`)
	rows, err := query.Query()
	if err != nil {
		return todos, err
	}

	defer rows.Close()
	for rows.Next() {
		var task data.Todo
		if err := rows.Scan(&task.Id, &task.Task, &task.Done); err != nil {
			return todos, err
		}

		todos = append(todos, task)
	}

	return todos, nil
}

func (this *PostgreSqlDataService) FindTodoByID(id int) (data.Todo, error) {
	var todo data.Todo
	query := this.queryService.Prepare(`SELECT * FROM "Todos" WHERE "Id" = $1`)
	if err := query.QueryRow(id).Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return todo, err
	}

	return todo, nil
}

func (this *PostgreSqlDataService) AddTodo(todo data.Todo) error {
	insert := this.queryService.Prepare(`INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2) RETURNING "Id"`)
	if _, err := insert.Exec(&todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

func (this *PostgreSqlDataService) UpdateTodo(todo data.Todo) error {
	update := this.queryService.Prepare(`UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`)
	if _, err := update.Exec(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

func (this *PostgreSqlDataService) RemoveTodo(id int) error {
	delete := this.queryService.Prepare(`DELETE FROM "Todos" WHERE "Id" = $1`)
	if _, err := delete.Exec(id); err != nil {
		return err
	}

	return nil

}

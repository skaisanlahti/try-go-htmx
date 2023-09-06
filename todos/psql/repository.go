package psql

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todos/domain"
)

type TodoRepository struct {
	Database *sql.DB
}

func NewTodoRepository(database *sql.DB) *TodoRepository {
	return &TodoRepository{database}
}

const selectTodos string = `SELECT * FROM "Todos" ORDER BY "Task" ASC`

func (repository *TodoRepository) GetTodos() []domain.Todo {
	var todos []domain.Todo
	rows, err := repository.Database.Query(selectTodos)
	if err != nil {
		log.Println(err.Error())
		return todos
	}

	defer rows.Close()
	for rows.Next() {
		var todo domain.Todo
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
			log.Println(err.Error())
			return todos
		}

		todos = append(todos, todo)
	}

	return todos
}

const selectTodoById string = `SELECT * FROM "Todos" WHERE "Id" = $1`

func (repository *TodoRepository) GetTodoById(id int) (domain.Todo, error) {
	var todo domain.Todo
	row := repository.Database.QueryRow(selectTodoById, id)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

const insertTodo string = `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2)`

func (repository *TodoRepository) AddTodo(todo domain.Todo) error {
	if _, err := repository.Database.Exec(insertTodo, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

const updateTodo string = `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`

func (repository *TodoRepository) UpdateTodo(todo domain.Todo) error {
	if _, err := repository.Database.Exec(updateTodo, &todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

const deleteTodo string = `DELETE FROM "Todos" WHERE "Id" = $1`

func (repository *TodoRepository) RemoveTodo(id int) error {
	if _, err := repository.Database.Exec(deleteTodo, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

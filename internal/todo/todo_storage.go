package todo

import (
	"database/sql"
	"log"
)

type todoStorage struct {
	database *sql.DB
}

func NewTodoStorage(db *sql.DB) *todoStorage {
	return &todoStorage{db}
}

func (storage *todoStorage) findTodos() []todo {
	var todos []todo
	query := `SELECT * FROM "Todos" ORDER BY "Task" ASC`
	rows, err := storage.database.Query(query)
	if err != nil {
		log.Println(err.Error())
		return todos
	}

	defer rows.Close()
	for rows.Next() {
		var todo todo
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
			log.Println(err.Error())
			return todos
		}

		todos = append(todos, todo)
	}

	return todos
}

func (storage *todoStorage) findTodoById(id int) (todo, error) {
	var todo todo
	query := `SELECT * FROM "Todos" WHERE "Id" = $1`
	row := storage.database.QueryRow(query, id)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func (storage *todoStorage) addTodo(todo todo) error {
	query := `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2)`
	if _, err := storage.database.Exec(query, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *todoStorage) updateTodo(todo todo) error {
	query := `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`
	if _, err := storage.database.Exec(query, &todo.Id, &todo.Task, &todo.Done); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *todoStorage) removeTodo(id int) error {
	query := `DELETE FROM "Todos" WHERE "Id" = $1`
	if _, err := storage.database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

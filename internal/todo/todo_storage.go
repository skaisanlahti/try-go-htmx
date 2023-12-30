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

func (this *todoStorage) listExists(name string) (entity.TodoList, error) {
	var todoList entity.TodoList
	query := `SELECT EXISTS(SELECT 1 FROM "TodoLists" WHERE "Name" = $1 )`
	row := this.database.QueryRow(query, name)
	if err := row.Scan(&todoList.Id, &todoList.Name, &todoList.UserId); err != nil {
		log.Println(err.Error())
		return todoList, err
	}

	return todoList, nil
}

func (this *todoStorage) findTodoListById(listId int) (entity.TodoList, error) {
	var todoList entity.TodoList
	query := `SELECT * FROM "TodoLists" WHERE "Id" = $1`
	row := this.database.QueryRow(query, listId)
	if err := row.Scan(&todoList.Id, &todoList.Name, &todoList.UserId); err != nil {
		log.Println(err.Error())
		return todoList, err
	}

	return todoList, nil
}

func (this *todoStorage) findTodoListsByUserId(userId int) []entity.TodoList {
	var todoLists []entity.TodoList
	query := `SELECT * FROM "TodoLists" WHERE "UserId" = $1 ORDER BY "Name" ASC`
	rows, err := this.database.Query(query, userId)
	if err != nil {
		log.Println(err.Error())
		return todoLists
	}

	defer rows.Close()
	for rows.Next() {
		var list entity.TodoList
		if err := rows.Scan(&list.Id, &list.Name, &list.UserId); err != nil {
			log.Println(err.Error())
			return todoLists
		}

		todoLists = append(todoLists, list)
	}

	return todoLists
}

func (this *todoStorage) insertTodoList(list entity.TodoList) error {
	query := `INSERT INTO "TodoLists" ("Name", "UserId") VALUES ($1, $2)`
	if _, err := this.database.Exec(query, list.Name, list.UserId); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (this *todoStorage) deleteTodoList(id int) error {
	query := `DELETE FROM "TodoLists" WHERE "Id" = $1`
	if _, err := this.database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (this *todoStorage) findTodosByListId(todoListId int) []entity.Todo {
	var todos []entity.Todo
	query := `SELECT * FROM "Todos" WHERE "TodoListId" = $1 ORDER BY "Task" ASC`
	rows, err := this.database.Query(query, todoListId)
	if err != nil {
		log.Println(err.Error())
		return todos
	}

	defer rows.Close()
	for rows.Next() {
		var todo entity.Todo
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done, &todo.TodoListId); err != nil {
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
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done, &todo.TodoListId); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func (this *todoStorage) insertTodo(todo entity.Todo) error {
	query := `INSERT INTO "Todos" ("Task", "Done", "TodoListId") VALUES ($1, $2, $3)`
	if _, err := this.database.Exec(query, todo.Task, todo.Done, todo.TodoListId); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (this *todoStorage) updateTodo(todo entity.Todo) error {
	query := `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`
	if _, err := this.database.Exec(query, todo.Id, todo.Task, todo.Done); err != nil {
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

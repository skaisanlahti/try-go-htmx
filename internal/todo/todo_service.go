package todo

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type TodoService struct {
	database *sql.DB
}

func NewTodoService(database *sql.DB) *TodoService {
	return &TodoService{database}
}

func (this *TodoService) FindListById(listId int) (entity.TodoList, error) {
	var todoList entity.TodoList
	query := `SELECT * FROM "TodoLists" WHERE "Id" = $1`
	row := this.database.QueryRow(query, listId)
	if err := row.Scan(&todoList.Id, &todoList.Name, &todoList.UserId); err != nil {
		log.Println(err.Error())
		return todoList, err
	}

	return todoList, nil
}

func (this *TodoService) FindListsByUserId(userId int) []entity.TodoList {
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

func (this *TodoService) FindTodosByListId(listId int) []entity.Todo {
	var todos []entity.Todo
	query := `SELECT * FROM "Todos" WHERE "TodoListId" = $1 ORDER BY "Task" ASC`
	rows, err := this.database.Query(query, listId)
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

func (this *TodoService) FindTodoById(id int) (entity.Todo, error) {
	var todo entity.Todo
	query := `SELECT * FROM "Todos" WHERE "Id" = $1`
	row := this.database.QueryRow(query, id)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done, &todo.TodoListId); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func (this *TodoService) AddList(name string, userId int) (entity.TodoList, error) {
	newList := entity.NewTodoList(name, userId)
	if err := newList.Validate(); err != nil {
		return newList, err
	}

	query := `INSERT INTO "TodoLists" ("Name", "UserId") VALUES ($1, $2)`
	if _, err := this.database.Exec(query, newList.Name, newList.UserId); err != nil {
		log.Println(err.Error())
		return newList, err
	}

	return newList, nil
}

func (this *TodoService) RemoveList(listId int) error {
	query := `DELETE FROM "TodoLists" WHERE "Id" = $1`
	if _, err := this.database.Exec(query, listId); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (this *TodoService) AddTodo(task string, todoListId int) (entity.Todo, error) {
	newTodo := entity.NewTodo(task, todoListId)
	if err := newTodo.Validate(); err != nil {
		return newTodo, err
	}

	query := `INSERT INTO "Todos" ("Task", "Done", "TodoListId") VALUES ($1, $2, $3)`
	if _, err := this.database.Exec(query, newTodo.Task, newTodo.Done, newTodo.TodoListId); err != nil {
		log.Println(err.Error())
		return newTodo, err
	}

	return newTodo, nil
}

func (this *TodoService) ToggleTodo(todoId int) (entity.Todo, error) {
	var todo entity.Todo
	tx, err := this.database.Begin()
	if err != nil {
		return todo, err
	}
	defer tx.Rollback()

	query := `SELECT * FROM "Todos" WHERE "Id" = $1 FOR UPDATE`
	row := tx.QueryRow(query, todoId)
	if err := row.Scan(&todo.Id, &todo.Task, &todo.Done, &todo.TodoListId); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	updatedTodo := todo.Toggle()
	query = `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`
	if _, err := tx.Exec(query, updatedTodo.Id, updatedTodo.Task, updatedTodo.Done); err != nil {
		log.Println(err.Error())
		return todo, err
	}

	if err = tx.Commit(); err != nil {
		return todo, err
	}

	return updatedTodo, nil
}

func (this *TodoService) RemoveTodo(todoId int) error {
	query := `DELETE FROM "Todos" WHERE "Id" = $1`
	if _, err := this.database.Exec(query, todoId); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

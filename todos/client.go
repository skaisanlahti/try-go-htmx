package todos

import (
	"database/sql"
	"sort"

	"github.com/skaisanlahti/test-go/common"
)

type todoClient interface {
	getAllTodos() ([]todoRecord, error)
	getTodoByID(id int) (todoRecord, error)
	addTodo(todo todoRecord) error
	updateTodo(todo todoRecord) error
	removeTodo(id int) error
}

type sqlTodoClient struct {
	queryClient *common.QueryClient
	database    *sql.DB
}

func newTodoDatabaseClient(db *sql.DB) *sqlTodoClient {
	return &sqlTodoClient{
		database:    db,
		queryClient: common.NewQueryClient(db),
	}
}

const getAllTodosQuery = `SELECT * FROM "Todos"`

func (this *sqlTodoClient) getAllTodos() ([]todoRecord, error) {
	var todos []todoRecord
	getAllTodos := this.queryClient.Prepare(getAllTodosQuery)
	rows, err := getAllTodos.Query()
	if err != nil {
		return todos, err
	}

	defer rows.Close()
	for rows.Next() {
		var todo todoRecord
		if err := rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	sort.SliceStable(todos, func(i int, j int) bool {
		return todos[i].Task < todos[j].Task
	})

	return todos, nil
}

const getTodoByIDQuery string = `SELECT * FROM "Todos" WHERE "Id" = $1`

func (this *sqlTodoClient) getTodoByID(id int) (todoRecord, error) {
	var todo todoRecord
	getTodoByID := this.queryClient.Prepare(getTodoByIDQuery)
	if err := getTodoByID.QueryRow(id).Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return todo, err
	}

	return todo, nil
}

const addTodoQuery string = `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2) RETURNING "Id"`

func (this *sqlTodoClient) addTodo(todo todoRecord) error {
	addTodo := this.queryClient.Prepare(addTodoQuery)
	if _, err := addTodo.Exec(&todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

const updateTodoQuery string = `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`

func (this *sqlTodoClient) updateTodo(todo todoRecord) error {
	updateTodo := this.queryClient.Prepare(updateTodoQuery)
	if _, err := updateTodo.Exec(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

const removeTodoQuery string = `DELETE FROM "Todos" WHERE "Id" = $1`

func (this *sqlTodoClient) removeTodo(id int) error {
	removeTodo := this.queryClient.Prepare(removeTodoQuery)
	if _, err := removeTodo.Exec(id); err != nil {
		return err
	}

	return nil
}

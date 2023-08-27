package todos

import (
	"database/sql"
	"sort"

	"github.com/skaisanlahti/test-go/common"
)

type repository[T any] interface {
	list() ([]T, error)
	find(id int) (T, error)
	add(record T) error
	update(record T) error
	remove(id int) error
}

type sqlTodoRepository struct {
	queryClient *common.QueryClient
	database    *sql.DB
}

func newSqlTodoRepository(db *sql.DB) *sqlTodoRepository {
	return &sqlTodoRepository{
		database:    db,
		queryClient: common.NewQueryClient(db),
	}
}

const getAllTodosQuery = `SELECT * FROM "Todos"`

func (this *sqlTodoRepository) list() ([]todoRecord, error) {
	var todos []todoRecord
	list := this.queryClient.Prepare(getAllTodosQuery)
	rows, err := list.Query()
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

func (this *sqlTodoRepository) find(id int) (todoRecord, error) {
	var todo todoRecord
	find := this.queryClient.Prepare(getTodoByIDQuery)
	if err := find.QueryRow(id).Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return todo, err
	}

	return todo, nil
}

const addTodoQuery string = `INSERT INTO "Todos" ("Task", "Done") VALUES ($1, $2) RETURNING "Id"`

func (this *sqlTodoRepository) add(todo todoRecord) error {
	add := this.queryClient.Prepare(addTodoQuery)
	if _, err := add.Exec(&todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

const updateTodoQuery string = `UPDATE "Todos" SET "Task" = $2, "Done" = $3 WHERE "Id" = $1`

func (this *sqlTodoRepository) update(todo todoRecord) error {
	update := this.queryClient.Prepare(updateTodoQuery)
	if _, err := update.Exec(&todo.Id, &todo.Task, &todo.Done); err != nil {
		return err
	}

	return nil
}

const removeTodoQuery string = `DELETE FROM "Todos" WHERE "Id" = $1`

func (this *sqlTodoRepository) remove(id int) error {
	remove := this.queryClient.Prepare(removeTodoQuery)
	if _, err := remove.Exec(id); err != nil {
		return err
	}

	return nil
}

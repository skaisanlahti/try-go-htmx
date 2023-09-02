package ports

import "database/sql"

type Query interface {
	Prepare(query string) *sql.Stmt
}

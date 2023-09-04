package interfaces

import "database/sql"

type QueryService interface {
	Prepare(query string) *sql.Stmt
}

package settings

import (
	"database/sql"
	"log"
)

func OpenDatabase(connectionString string) *sql.DB {
	database, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return database
}

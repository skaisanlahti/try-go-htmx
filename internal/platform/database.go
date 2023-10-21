package platform

import (
	"database/sql"
	"log"
	"os"
	"sort"
)

type DatabaseOptions struct {
	Driver             string
	ConnectionString   string
	MigrationDirectory string
	MigrateOnStartup   bool
}

func OpenDatabase(options DatabaseOptions) *sql.DB {
	database, err := sql.Open(options.Driver, options.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	if options.MigrateOnStartup {
		applyMigrations(database, options.MigrationDirectory)
	}

	return database
}

func applyMigrations(database *sql.DB, directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err.Error())
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		content, err := os.ReadFile(directory + "/" + file.Name())
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = database.Exec(string(content))
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Executed migration: %s", file.Name())
	}
}

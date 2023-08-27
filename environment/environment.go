package environment

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

type Environment struct {
	Address  string
	Database string
}

func Read(filename string) Environment {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Print("Env file not found")
		log.Fatal(err)
	}

	envMap := make(map[string]string)
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envMap[key] = value
	}

	address, exists := envMap["ADDRESS"]
	if !exists {
		log.Fatal("ADDRESS not found in env file.")
	}

	database, exists := envMap["DATABASE"]
	if !exists {
		log.Fatal("DATABASE not found in env file.")
	}

	return Environment{
		Address:  address,
		Database: database,
	}
}

func OpenDatabase(connectionString string) *sql.DB {
	database, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return database
}

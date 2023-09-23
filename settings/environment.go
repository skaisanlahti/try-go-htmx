package settings

import (
	"log"
	"os"
	"strings"
)

const (
	ModeDevelopment string = "Development"
	ModeProduction  string = "Production"
)

func IsDevelopment(mode string) bool {
	return mode == ModeDevelopment
}

func IsProduction(mode string) bool {
	return mode == ModeProduction
}

type Environment struct {
	Mode     string
	Address  string
	Database string
}

func ReadEnvironment(filename string) Environment {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Print("Env file not found")
		log.Fatal(err)
	}

	variables := make(map[string]string)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		variables[key] = value
	}

	mode, exists := variables["MODE"]
	if !exists {
		log.Println("MODE not found in env file. Using Development.")
		mode = ModeDevelopment
	}

	address, exists := variables["ADDRESS"]
	if !exists {
		log.Fatal("ADDRESS not found in env file.")
	}

	database, exists := variables["DATABASE"]
	if !exists {
		log.Fatal("DATABASE not found in env file.")
	}

	return Environment{
		Mode:     mode,
		Address:  address,
		Database: database,
	}
}

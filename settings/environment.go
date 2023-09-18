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

	mode, exists := envMap["MODE"]
	if !exists {
		log.Fatal("MODE not found in env file.")
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
		Mode:     mode,
		Address:  address,
		Database: database,
	}
}

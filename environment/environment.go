package environment

import (
	"log"
	"os"
	"strings"
)

type Environment struct {
	Address string
}

func New(filename string) Environment {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Env file not found in")
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

	return Environment{
		Address: address,
	}
}

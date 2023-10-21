package platform

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type settings struct {
	Mode     string           `json:"mode"`
	Address  string           `json:"address"`
	Database databaseSettings `json:"database"`
	Password passwordSettings `json:"password"`
	Session  sessionSettings  `json:"session"`
}

type databaseSettings struct {
	ConnectionString   string `json:"connectionString"`
	MigrationDirectory string `json:"migrationDirectory"`
	MigrateOnStartup   bool   `json:"migrateOnStartup"`
}

type passwordSettings struct {
	Cost                int    `json:"cost"`
	Time                uint32 `json:"time"`
	Memory              uint32 `json:"memory"`
	Threads             uint8  `json:"threads"`
	SaltLength          uint32 `json:"saltLength"`
	KeyLength           uint32 `json:"keyLength"`
	RecalculateOutdated bool   `json:"recalculateOutdated"`
}

type sessionSettings struct {
	Secure             bool    `json:"secure"`
	CookieName         string  `json:"cookieName"`
	SecretLength       uint32  `json:"secretLength"`
	SessionDurationMin float64 `json:"sessionDurationMin"`
}

func ReadSettings(file string) settings {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err.Error())
	}

	var settings settings
	if err := json.Unmarshal(bytes, &settings); err != nil {
		log.Panic(err.Error())
	}

	return settings
}

func IsDevelopment(mode string) bool {
	return strings.ToLower(mode) == "development"
}

func IsProduction(mode string) bool {
	return strings.ToLower(mode) == "production"
}

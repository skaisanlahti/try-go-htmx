package todoapp

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Settings struct {
	Mode     string           `json:"mode"`
	Address  string           `json:"address"`
	Database DatabaseSettings `json:"database"`
	Password PasswordSettings `json:"password"`
	Session  SessionSettings  `json:"session"`
}

type DatabaseSettings struct {
	ConnectionString   string `json:"connectionString"`
	MigrationDirectory string `json:"migrationDirectory"`
	MigrateOnStartup   bool   `json:"migrateOnStartup"`
}

type PasswordSettings struct {
	Cost                int    `json:"cost"`
	Time                uint32 `json:"time"`
	Memory              uint32 `json:"memory"`
	Threads             uint8  `json:"threads"`
	SaltLength          uint32 `json:"saltLength"`
	KeyLength           uint32 `json:"keyLength"`
	RecalculateOutdated bool   `json:"recalculateOutdated"`
}

type SessionSettings struct {
	Secure             bool    `json:"secure"`
	CookieName         string  `json:"cookieName"`
	SecretLength       uint32  `json:"secretLength"`
	SessionDurationMin float64 `json:"sessionDurationMin"`
}

func ReadSettings(file string) Settings {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err.Error())
	}

	var settings Settings
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

package todoapp

import (
	"log"
	"time"
)

type PasswordService interface {
	CreateKey(password string) ([]byte, error)
	VerifyKey(encodedKey []byte, candidatePassword string) (bool, chan []byte)
}

const (
	minimumDurationMs int64 = 200
	maximumDurationMs int64 = 500
)

func MonitorEncodingTime() func() {
	start := time.Now()
	return func() {
		durationMs := time.Now().Sub(start).Milliseconds()
		if durationMs < minimumDurationMs {
			log.Printf("Password encoding took less than %d ms (%d ms). Consider increasing encoding difficult.", minimumDurationMs, durationMs)
		}

		if durationMs > maximumDurationMs {
			log.Printf("Password encoding took more than %d ms (%d ms). Consider decreasing encoding difficult.", maximumDurationMs, durationMs)
		}
	}
}

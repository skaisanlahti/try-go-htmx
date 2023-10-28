package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id      string
	UserId  int
	Expires time.Time
}

func NewSession(userId int, duration time.Duration) Session {
	sessionId := uuid.New().String()
	expires := time.Now().Add(duration)
	return Session{sessionId, userId, expires}
}

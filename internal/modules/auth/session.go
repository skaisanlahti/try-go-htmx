package auth

import (
	"time"

	"github.com/google/uuid"
)

type session struct {
	Id      string
	UserId  int
	Expires time.Time
}

func newSession(userId int, duration time.Duration) session {
	sessionId := uuid.New().String()
	expires := time.Now().Add(duration)
	return session{sessionId, userId, expires}
}

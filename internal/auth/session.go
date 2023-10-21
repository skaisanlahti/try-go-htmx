package auth

import (
	"time"
)

type session struct {
	Id      string
	UserId  int
	Expires time.Time
}

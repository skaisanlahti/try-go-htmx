package mem

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type SessionAccessor struct {
	sessions map[string]*todoapp.Session
	locker   sync.RWMutex
}

func NewSessionAccessor() *SessionAccessor {
	accessor := &SessionAccessor{
		sessions: make(map[string]*todoapp.Session),
	}

	go removeExpired(accessor)
	return accessor
}

const (
	checkingInterval time.Duration = 60 * time.Second
	timeFormat       string        = "2006/01/02 15:04:05 -0700"
)

func removeExpired(accessor *SessionAccessor) {
	log.Printf("Started a session clean up process at %s.", time.Now().Format(timeFormat))
	for {
		log.Printf("Next expired session clean up scheduled at %s.", time.Now().Add(checkingInterval).Format(timeFormat))
		time.Sleep(checkingInterval)
		startTask := time.Now()
		accessor.locker.Lock()
		for _, session := range accessor.sessions {
			if session.Expires.Before(time.Now()) {
				delete(accessor.sessions, session.Id)
			}
		}

		accessor.locker.Unlock()
		taskDuration := time.Now().Sub(startTask)
		log.Printf("Expired sessions cleaned up in %d ms.", taskDuration.Milliseconds())
	}
}

func (accessor *SessionAccessor) FindSession(sessionId string) (*todoapp.Session, error) {
	accessor.locker.RLock()
	defer accessor.locker.RUnlock()
	session, ok := accessor.sessions[sessionId]
	if !ok {
		return nil, errors.New("Session not found.")
	}

	return session, nil
}

func (accessor *SessionAccessor) AddSession(newSession *todoapp.Session) error {
	accessor.locker.Lock()
	defer accessor.locker.Unlock()
	for _, session := range accessor.sessions {
		if session.UserId == newSession.UserId {
			delete(accessor.sessions, newSession.Id)
		}
	}

	accessor.sessions[newSession.Id] = newSession
	return nil
}

func (accessor *SessionAccessor) UpdateSession(session *todoapp.Session) error {
	accessor.locker.Lock()
	defer accessor.locker.Unlock()
	accessor.sessions[session.Id] = session
	return nil
}

func (accessor *SessionAccessor) RemoveSession(sessionId string) error {
	accessor.locker.Lock()
	defer accessor.locker.Unlock()
	delete(accessor.sessions, sessionId)
	return nil
}

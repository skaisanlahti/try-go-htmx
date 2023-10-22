package auth

import (
	"errors"
	"log"
	"sync"
	"time"
)

type sessionStorage struct {
	sessions map[string]session
	locker   sync.RWMutex
}

func newSessionStorage() *sessionStorage {
	storage := &sessionStorage{
		sessions: make(map[string]session),
	}

	go storage.removeExpired()
	return storage
}

func (storage *sessionStorage) findSession(sessionId string) (session, error) {
	storage.locker.RLock()
	defer storage.locker.RUnlock()
	session, ok := storage.sessions[sessionId]
	if !ok {
		return session, errors.New("Session not found.")
	}

	return session, nil
}

func (storage *sessionStorage) addSession(newSession session) error {
	storage.locker.Lock()
	defer storage.locker.Unlock()
	for _, session := range storage.sessions {
		if session.UserId == newSession.UserId {
			delete(storage.sessions, newSession.Id)
		}
	}

	storage.sessions[newSession.Id] = newSession
	return nil
}

func (storage *sessionStorage) updateSession(session session) error {
	storage.locker.Lock()
	defer storage.locker.Unlock()
	storage.sessions[session.Id] = session
	return nil
}

func (storage *sessionStorage) removeSession(sessionId string) error {
	storage.locker.Lock()
	defer storage.locker.Unlock()
	delete(storage.sessions, sessionId)
	return nil
}

const (
	checkingInterval time.Duration = 60 * time.Second
	timeFormat       string        = "2006/01/02 15:04:05 -0700"
)

func (storage *sessionStorage) removeExpired() {
	log.Printf("Started a session clean up process at %s.", time.Now().Format(timeFormat))
	for {
		log.Printf("Next expired session clean up scheduled at %s.", time.Now().Add(checkingInterval).Format(timeFormat))
		time.Sleep(checkingInterval)
		startTask := time.Now()
		storage.locker.Lock()
		for _, session := range storage.sessions {
			if session.Expires.Before(time.Now()) {
				delete(storage.sessions, session.Id)
			}
		}

		storage.locker.Unlock()
		taskDuration := time.Now().Sub(startTask)
		log.Printf("Expired sessions cleaned up in %d ms.", taskDuration.Milliseconds())
	}
}

package security

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type sessionStorage struct {
	sessions map[string]*entity.Session
	locker   sync.RWMutex
}

func newSessionStorage() *sessionStorage {
	storage := &sessionStorage{
		sessions: make(map[string]*entity.Session),
	}

	go storage.removeExpired()
	return storage
}

func (storage *sessionStorage) findSessionBySessionId(sessionId string) (entity.Session, error) {
	storage.locker.RLock()
	defer storage.locker.RUnlock()
	session, ok := storage.sessions[sessionId]
	if !ok {
		return *session, errors.New("Session not found.")
	}

	return *session, nil
}

func (storage *sessionStorage) findSessionByUserId(userId int) (entity.Session, error) {
	storage.locker.RLock()
	defer storage.locker.RUnlock()
	session, ok := storage.sessions[strconv.Itoa(userId)]
	if !ok {
		return *session, errors.New("Session not found.")
	}

	return *session, nil
}

func (storage *sessionStorage) insertSession(newSession entity.Session) error {
	storage.locker.Lock()
	defer storage.locker.Unlock()
	session, ok := storage.sessions[newSession.Id]
	if ok {
		delete(storage.sessions, strconv.Itoa(session.UserId))
		delete(storage.sessions, session.Id)
	}

	storage.sessions[newSession.Id] = &newSession
	storage.sessions[strconv.Itoa(newSession.UserId)] = &newSession
	return nil
}

func (storage *sessionStorage) updateSession(session entity.Session) error {
	storage.locker.Lock()
	defer storage.locker.Unlock()
	storage.sessions[session.Id] = &session
	storage.sessions[strconv.Itoa(session.UserId)] = &session
	return nil
}

func (storage *sessionStorage) deleteSession(session entity.Session) error {
	storage.locker.Lock()
	defer storage.locker.Unlock()
	delete(storage.sessions, session.Id)
	delete(storage.sessions, strconv.Itoa(session.UserId))
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
				delete(storage.sessions, strconv.Itoa(session.UserId))
			}
		}

		storage.locker.Unlock()
		taskDuration := time.Now().Sub(startTask)
		log.Printf("Expired sessions cleaned up in %d ms.", taskDuration.Milliseconds())
	}
}

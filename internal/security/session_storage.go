package security

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type SessionOptions struct {
	Secure     bool
	CookieName string
	Secret     string
	Duration   time.Duration
}

type SessionStorage struct {
	sessions map[string]entity.Session
	locker   sync.RWMutex
}

func NewSessionStorage() *SessionStorage {
	storage := &SessionStorage{
		sessions: make(map[string]entity.Session),
	}

	go storage.RemoveExpired()
	return storage
}

func (this *SessionStorage) FindSessionBySessionId(sessionId string) (entity.Session, error) {
	this.locker.RLock()
	defer this.locker.RUnlock()
	session, ok := this.sessions[sessionId]
	if !ok {
		return session, errors.New("Session not found.")
	}

	return session, nil
}

func (this *SessionStorage) FindSessionByUserId(userId int) (entity.Session, error) {
	this.locker.RLock()
	defer this.locker.RUnlock()
	session, ok := this.sessions[strconv.Itoa(userId)]
	if !ok {
		return session, errors.New("Session not found.")
	}

	return session, nil
}

func (this *SessionStorage) InsertSession(newSession entity.Session) error {
	this.locker.Lock()
	defer this.locker.Unlock()
	session, ok := this.sessions[newSession.Id]
	if ok {
		delete(this.sessions, strconv.Itoa(session.UserId))
		delete(this.sessions, session.Id)
	}

	this.sessions[newSession.Id] = newSession
	this.sessions[strconv.Itoa(newSession.UserId)] = newSession
	return nil
}

func (this *SessionStorage) UpdateSession(session entity.Session) error {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.sessions[session.Id] = session
	this.sessions[strconv.Itoa(session.UserId)] = session
	return nil
}

func (this *SessionStorage) DeleteSession(sessionId string) error {
	this.locker.Lock()
	defer this.locker.Unlock()

	session, ok := this.sessions[sessionId]
	if !ok {
		return errors.New("Session not found.")
	}

	delete(this.sessions, session.Id)
	delete(this.sessions, strconv.Itoa(session.UserId))
	return nil
}

const (
	checkingInterval time.Duration = 60 * time.Second
	timeFormat       string        = "2006/01/02 15:04:05 -0700"
)

func (this *SessionStorage) RemoveExpired() {
	log.Printf("Started a session clean up process at %s.", time.Now().Format(timeFormat))
	for {
		log.Printf("Next expired session clean up scheduled at %s.", time.Now().Add(checkingInterval).Format(timeFormat))
		time.Sleep(checkingInterval)
		startTask := time.Now()
		this.locker.Lock()
		for _, session := range this.sessions {
			if session.Expires.Before(time.Now()) {
				delete(this.sessions, session.Id)
				delete(this.sessions, strconv.Itoa(session.UserId))
			}
		}

		this.locker.Unlock()
		taskDuration := time.Now().Sub(startTask)
		log.Printf("Expired sessions cleaned up in %d ms.", taskDuration.Milliseconds())
	}
}

package sessions

import (
	"errors"
	"log"
	"sync"
	"time"
)

type InMemoryRepository struct {
	sessions map[string]*Session
	locker   sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{
		sessions: make(map[string]*Session),
	}

	go removeExpired(repository)
	return repository
}

func (repository *InMemoryRepository) Find(sessionId string) (*Session, error) {
	repository.locker.RLock()
	defer repository.locker.RUnlock()
	session, ok := repository.sessions[sessionId]
	if !ok {
		return nil, errors.New("Session not found.")
	}

	return session, nil
}

func (repository *InMemoryRepository) Add(session *Session) error {
	repository.locker.Lock()
	defer repository.locker.Unlock()
	for _, s := range repository.sessions {
		if s.UserId == session.UserId {
			delete(repository.sessions, session.Id)
		}
	}

	repository.sessions[session.Id] = session
	return nil
}

func (repository *InMemoryRepository) Update(session *Session) error {
	repository.locker.Lock()
	defer repository.locker.Unlock()
	repository.sessions[session.Id] = session
	return nil
}

func (repository *InMemoryRepository) Remove(sessionId string) error {
	repository.locker.Lock()
	defer repository.locker.Unlock()
	delete(repository.sessions, sessionId)
	return nil
}

const (
	checkingInterval time.Duration = 60 * time.Second // seconds
	timeFormat       string        = "2006-01-02 15:04:05 MST"
)

func removeExpired(repository *InMemoryRepository) {
	log.Printf("Started a session clean up process at %s.", time.Now().Format(timeFormat))
	for {
		log.Printf("Next expired session clean up scheduled at %s.", time.Now().Add(checkingInterval).Format(timeFormat))
		time.Sleep(checkingInterval)
		startTask := time.Now()
		repository.locker.Lock()
		for _, session := range repository.sessions {
			if session.Expires.Before(time.Now()) {
				delete(repository.sessions, session.Id)
			}
		}

		repository.locker.Unlock()
		taskDuration := time.Now().Sub(startTask)
		log.Printf("Expired sessions cleaned up in %d ms.", taskDuration.Milliseconds())
	}
}

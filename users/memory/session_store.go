package memory

import (
	"log"
	"sync"
	"time"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type SessionStore struct {
	sessions map[string]*domain.Session
	locker   sync.Mutex
}

func NewSessionStore() *SessionStore {
	store := &SessionStore{sessions: make(map[string]*domain.Session)}
	go cleanExpiredSessions(store)
	return store
}

func (store *SessionStore) Add(userId int) *domain.Session {
	store.locker.Lock()
	defer store.locker.Unlock()
	var existingSession *domain.Session
	for _, session := range store.sessions {
		if session.UserId == userId {
			existingSession = session
			break
		}
	}

	if existingSession != nil {
		store.Remove(existingSession.Id)
	}

	newSession := domain.NewSession(userId)
	store.sessions[newSession.Id] = newSession
	return newSession
}

func (store *SessionStore) Remove(sessionId string) {
	store.locker.Lock()
	defer store.locker.Unlock()
	delete(store.sessions, sessionId)
}

func (store *SessionStore) Validate(sessionId string) (*domain.Session, bool) {
	store.locker.Lock()
	defer store.locker.Unlock()
	session, ok := store.sessions[sessionId]
	if !ok {
		return nil, false
	}

	if session.Expires.Before(time.Now()) {
		store.Remove(session.Id)
		return nil, false
	}

	return session, true
}

func (store *SessionStore) Extend(session *domain.Session) *domain.Session {
	store.locker.Lock()
	defer store.locker.Unlock()
	session.Expires = time.Now().Add(domain.SessionDuration)
	store.sessions[session.Id] = session
	return session
}

const (
	checkingInterval time.Duration = 60 * time.Second // seconds
	timeFormat       string        = "2006-01-02 15:04:05 MST"
)

func cleanExpiredSessions(store *SessionStore) {
	log.Printf("Started a session clean up process at %s.", time.Now().Format(timeFormat))
	for {
		log.Printf("Next expired session clean up scheduled at %s.", time.Now().Add(checkingInterval).Format(timeFormat))
		time.Sleep(checkingInterval)
		startTask := time.Now()
		store.locker.Lock()
		for _, session := range store.sessions {
			if session.Expires.Before(time.Now()) {
				store.Remove(session.Id)
			}
		}

		store.locker.Unlock()
		taskDuration := time.Now().Sub(startTask)
		log.Printf("Expired sessions cleaned up in %d ms.", taskDuration.Milliseconds())
	}
}

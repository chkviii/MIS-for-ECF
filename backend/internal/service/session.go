package service

/*
keep session info in memory temporarily,
this is not a persistent storage solution.
*/

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string `json:"session_id"`
	Key       []byte `json:"session_key"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

var Sess *sync.Map // key: string, value: *Session

func InitSessionService() {
	// Initialize the session service with an empty map
	Sess = &sync.Map{}
	go func() {
		for {
			time.Sleep(5 * time.Minute) // Check every 5 minutes
			Sess.Range(func(key, value any) bool {
				session, ok := value.(*Session)
				if !ok || session == nil {
					return true // Skip if not a valid session
				}
				if !session.isAlive() {
					Sess.Delete(key) // Remove expired session
				}
				return true // Continue iteration
			})
		}
	}()
}

func NewSession(sessions *sync.Map, UserID string) *Session {
	// Generate a new session ID and key
	// check if session already exists
	newID := uuid.New().String()
	for _, ok := sessions.Load(newID); ok; {
		newID = uuid.New().String()
	}
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		panic("failed to generate session key")
	}

	newSession := &Session{
		ID:        newID,
		Key:       newKey,
		UserID:    UserID,
		ExpiresAt: time.Now().Add(60 * time.Minute).Unix(), // 60 minutes expiration
	}
	sessions.Store(newID, newSession)

	return newSession
}

func (s *Session) isAlive() bool {
	// Check if the session is still valid based on ExpiresAt
	return s.ExpiresAt > time.Now().Unix()
}

func (s *Session) Refresh() {
	// Refresh the session expiration time
	s.ExpiresAt = time.Now().Add(60 * time.Minute).Unix() // 60 minutes expiration
}

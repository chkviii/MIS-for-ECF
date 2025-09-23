package service

/*
keep session info in memory temporarily,
this is not a persistent storage solution.
*/

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         string `json:"session_id"`
	Key        []byte `json:"session_key"`
	UserID     string `json:"user_id"`
	UserPubKey string `json:"user_pub_key,omitempty"`
	ExpiresAt  int64  `json:"expires_at"`
}

type SessionManager struct {
	sess *sync.Map
}

var (
	SessMgr  *SessionManager // Global instance of SessionManager
	initOnce sync.Once
)

func GetSessMgr() *SessionManager {
	InitSessionService()
	return SessMgr
}

func InitSessionService() {
	initOnce.Do(func() {
		SessMgr = &SessionManager{
			sess: &sync.Map{},
		}
		go SessMgr.startCleanupRoutine()
	})
}

func (sm *SessionManager) startCleanupRoutine() {
	for {
		time.Sleep(5 * time.Minute) // Cleanup interval
		sm.sess.Range(func(key, value any) bool {
			session, ok := value.(*Session)
			if !ok || session == nil {
				return true
			}
			if !session.isAlive() {
				sm.sess.Delete(key)
			}
			return true
		})
	}
}

func (sm *SessionManager) NewSession(UserID string) *Session {
	// Generate a new session ID and key
	// check if session already exists
	newID := uuid.New().String()
	for _, ok := sm.sess.Load(newID); !ok; newID = uuid.New().String() {
	}

	// Generate a random 32-byte key for AES-256
	newKey := make([]byte, 32)
	rand.Read(newKey) // ignore error, unlikely to happen
	// keyB64 := base64.StdEncoding.EncodeToString(newKey)

	newSession := &Session{
		ID:        newID,
		Key:       newKey,
		UserID:    UserID,
		ExpiresAt: time.Now().Add(2 * time.Minute).Unix(), // 2 minutes expiration
	}

	sm.sess.Store(newID, newSession)

	return newSession
}

func (sm *SessionManager) NewRegSession(UserID string) *Session {
	// Generate a new session ID and key
	// check if session already exists
	newID := uuid.New().String()
	for _, ok := sm.sess.Load(newID); ok; newID = uuid.New().String() {
	}

	// Generate a random 32-byte key for AES-256
	newKey := make([]byte, 32)
	rand.Read(newKey) // ignore error, unlikely to happen
	// keyB64 := base64.StdEncoding.EncodeToString(newKey)

	newSession := &Session{
		ID:        newID,
		Key:       newKey,
		UserID:    UserID,
		ExpiresAt: time.Now().Add(60 * time.Minute).Unix(), // 60 minutes expiration
	}

	sm.sess.Store(newID, newSession)

	return newSession
}

func (sm *SessionManager) getSession(sessionID string) (*Session, error) {
	value, ok := sm.sess.Load(sessionID)
	if !ok {
		return nil, errors.New("session not found")
	}
	session, ok := value.(*Session)
	if !ok || session == nil || !session.isAlive() {
		sm.sess.Delete(sessionID)
		return nil, errors.New("session expired or invalid")
	}
	return session, nil
}

func (s *Session) isAlive() bool {
	// Check if the session is still valid based on ExpiresAt
	return s.ExpiresAt > time.Now().Unix()
}

func (s *Session) Refresh() {
	// Refresh the session expiration time
	if s.ExpiresAt < time.Now().Add(60*time.Minute).Unix() { // only extend if less than 60 minutes left
		s.ExpiresAt = time.Now().Add(60 * time.Minute).Unix() // 60 minutes expiration
	}
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.sess.Delete(sessionID)
}

// Not required now
func (sm *SessionManager) GetUser(sessionID string) (string, error) {
	session, err := sm.getSession(sessionID)
	if err != nil {
		return "", err
	}
	return session.UserID, nil
}

func (sm *SessionManager) GetKey(sessionID string) ([]byte, error) {

	session, err := sm.getSession(sessionID)
	if err != nil {
		return nil, err
	}

	return session.Key, nil
}

func (sm *SessionManager) GetExpTime(sessionID string) (int64, error) {
	session, err := sm.getSession(sessionID)
	if err != nil {
		return -1, err
	}
	return session.ExpiresAt, nil
}

func (sm *SessionManager) Refresh(sessionID string, minutes int) (int64, error) {
	session, err := sm.getSession(sessionID)
	if err != nil {
		return -1, err
	}
	t := time.Now().Add(time.Duration(minutes) * time.Minute).Unix()

	if session.ExpiresAt < t {
		session.ExpiresAt = t
	}

	return session.ExpiresAt, nil
}

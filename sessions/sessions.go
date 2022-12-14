package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

func GenerateSessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

type Session struct {
	sessionId       string                 // session id唯一标示
	recentlyVisited time.Time              // 最后访问时间
	value           map[string]interface{} // session里面存储的值
}

func (session *Session) Set(key string, value interface{}) error {
	session.value[key] = value

	return nil
}

func (session *Session) Get(key string) interface{} {

	if v, ok := session.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (session *Session) Delete(key string) error {
	delete(session.value, key)

	return nil
}

func (session *Session) SessionID() string {
	return session.sessionId
}

type MemoryStore struct {
	lock     sync.Mutex
	sessions map[string]*Session
}

func InitMemoryStore() *MemoryStore {
	sessionss := make(map[string]*Session, 0)
	return &MemoryStore{sessions: sessionss}
}

func (ms *MemoryStore) Add(sessionId string, Value map[string]interface{}) (*Session, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	newsess := &Session{sessionId: sessionId, recentlyVisited: time.Now(), value: Value}
	ms.sessions[sessionId] = newsess
	return newsess, nil
}

func (ms *MemoryStore) Get(sessionId string) (*Session, error) {
	if element, ok := ms.sessions[sessionId]; ok {
		return element, nil
	} else {
		return nil, errors.New("session not found in RAM")
	}
}

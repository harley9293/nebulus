package httpd

import (
	"fmt"
	"time"
)

type Session struct {
	id     string
	values map[string]any
}

func (s *Session) Get(key string) any {
	return s.values[key]
}

func (s *Session) Set(key string, value any) {
	s.values[key] = value
}

type sessionMng struct {
	data map[string]*Session
}

func newSessionMng() *sessionMng {
	return &sessionMng{
		data: make(map[string]*Session),
	}
}

func (m *sessionMng) get(id string) *Session {
	if session, ok := m.data[id]; ok {
		return session
	} else {
		return nil
	}
}

func (m *sessionMng) new(key string) *Session {
	session := &Session{
		id:     generateSessionID(key),
		values: make(map[string]any),
	}
	m.data[session.id] = session
	return session
}

func generateSessionID(key string) string {

	return fmt.Sprintf("%s%d", key, time.Now().UnixNano())
}

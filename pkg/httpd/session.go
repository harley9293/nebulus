package httpd

import (
	"fmt"
	"github.com/harley9293/nebulus/pkg/db"
	"time"
)

type Session interface {
	Get(key string) string
	Set(key, value string)
	UpdateExpire()
}

type localSession struct {
	data       map[string]string
	expireTime time.Time

	cfgExpireTime time.Duration
}

func (s *localSession) Get(key string) string {
	return s.data[key]
}

func (s *localSession) Set(key, value string) {
	s.data[key] = value
}

func (s *localSession) UpdateExpire() {
	s.expireTime = time.Now().Add(s.cfgExpireTime)
}

type redisSession struct {
	id string
	*db.RedisClient

	cfgExpireTime time.Duration
}

func (s *redisSession) Get(key string) string {
	if value, err := s.HGet(s.id, key).Result(); err != nil {
		return ""
	} else {
		return value
	}
}

func (s *redisSession) Set(key, value string) {
	s.HSet(s.id, key, value)
}

func (s *redisSession) UpdateExpire() {
	s.Expire(s.id, s.cfgExpireTime)
}

type sessionMng struct {
	session map[string]Session

	sType      SessionType
	expireTime time.Duration
	redis      *db.RedisConfig
}

func newSessionMng(sType SessionType, expireTime time.Duration, redis *db.RedisConfig) *sessionMng {
	return &sessionMng{
		session:    make(map[string]Session),
		sType:      sType,
		expireTime: expireTime,
		redis:      redis,
	}
}

func (m *sessionMng) get(id string) Session {
	if session, ok := m.session[id]; ok {
		return session
	} else {
		return nil
	}
}

func (m *sessionMng) new(key string) Session {
	id := generateSessionID(key)
	switch m.sType {
	case SessionTypeLocal:
		m.session[id] = &localSession{
			data:          make(map[string]string),
			expireTime:    time.Now().Add(m.expireTime),
			cfgExpireTime: m.expireTime,
		}
	case SessionTypeRedis:
		m.session[id] = &redisSession{
			id:            id,
			RedisClient:   db.NewRedisClient(m.redis),
			cfgExpireTime: m.expireTime,
		}
		m.session[id].UpdateExpire()
	}

	m.session[id].Set("id", id)
	return m.session[id]
}

func generateSessionID(key string) string {
	return fmt.Sprintf("%s%d", key, time.Now().UnixNano())
}

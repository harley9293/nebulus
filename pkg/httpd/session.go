package httpd

import (
	"fmt"
	"github.com/harley9293/nebulus/pkg/def"
	"time"
)

type defaultSession struct {
	id         string
	data       map[string]any
	expireTime time.Time

	cfgExpireTime time.Duration
}

func (s *defaultSession) New(key string) def.Session {
	return &defaultSession{
		id:         fmt.Sprintf("%s%d", key, time.Now().UnixNano()),
		data:       make(map[string]any),
		expireTime: time.Now().Add(s.cfgExpireTime),

		cfgExpireTime: s.cfgExpireTime,
	}
}

func (s *defaultSession) ID() string {
	return s.id
}

func (s *defaultSession) Get(key string) any {
	return s.data[key]
}

func (s *defaultSession) Set(key string, value any) {
	s.data[key] = value
}

func (s *defaultSession) UpdateExpire() {
	s.expireTime = time.Now().Add(s.cfgExpireTime)
}

func (s *defaultSession) IsExpired() bool {
	return s.expireTime.Before(time.Now())
}

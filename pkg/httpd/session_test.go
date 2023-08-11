package httpd

import (
	"testing"
	"time"
)

func TestDefaultSession(t *testing.T) {
	b := &defaultSession{cfgExpireTime: 24 * time.Hour}
	s := b.New("test")
	if s.Get("test") != nil {
		t.Fatal("Get() failed")
	}

	s.Set("test", "test")
	if s.Get("test") != "test" {
		t.Fatal("Get() failed")
	}

	s.UpdateExpire()
}

package httpd

import "testing"

func TestSession(t *testing.T) {
	mng := newSessionMng(SessionTypeLocal, 24*3600, nil)
	session := mng.new("hello")
	if mng.get(session.Get("id")) == nil {
		t.Fatal("session not found")
	}

	if mng.get("world") != nil {
		t.Fatal("session found")
	}

	session.Set("test", "one")
	if session.Get("test") != "one" {
		t.Fatal("session set/get error")
	}

	if session.Get("test2") != "" {
		t.Fatal("session set/get error")
	}
}

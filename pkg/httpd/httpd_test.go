package httpd

import (
	"github.com/harley9293/nebulus"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go nebulus.Run()
	time.Sleep(10 * time.Millisecond)
	m.Run()
	nebulus.Shutdown()
}

func NewTestHttpService(name string, host string, f func(s *Service)) {
	s := NewService()
	f(s)
	err := nebulus.Register(name, s, host)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Millisecond)
}

type EmptyTestStruct struct {
}

func TestService_AddHandler_Fail(t *testing.T) {
	defer func() { recover() }()

	service := NewService()
	service.AddHandler("GET", "/test", 1)
}

func TestService_OnInit(t *testing.T) {
	service := NewService()
	err := service.OnInit()
	if err == nil {
		t.Fatal("OnInit() failed, err is nil")
	}

	err = service.OnInit(123456)
	if err == nil {
		t.Fatal("OnInit() failed, err is nil")
	}

	err = service.OnInit("http://localhost:80")
	if err == nil {
		t.Fatal("OnInit() failed, err:" + err.Error())
	}
}

func TestServiceFailed(t *testing.T) {
	s := NewService()
	s.AddHandler("GET", "/test", func(req *EmptyTestStruct, ctx *Context) string {
		return "hello world"
	})
	err := nebulus.Register("Test", s, "127.0.0.1:30011")
	if err != nil {
		t.Fatal(err)
	}
	defer nebulus.Destroy("Test")
	time.Sleep(10 * time.Millisecond)
	err = s.srv.Close()
	if err != nil {
		t.Fatal("srv.Close() failed, err:" + err.Error())
	}
	time.Sleep(1 * time.Second)

	client := NewClient("http://127.0.0.1:30011")
	err = client.Get("/test", nil)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}
	if client.strRsp != "hello world" {
		t.Fatal("rsp not ok, rsp:" + client.strRsp)
	}
}

func TestService_UseSession(t *testing.T) {
	s := NewService()
	s.UseSession(&defaultSession{})

	s2 := s.NewSession("test")
	if s.GetSession(s2.ID()) != nil {
		t.Fatal("session not nil")
	}
	if s.GetSession("test") != nil {
		t.Fatal("session not nil")
	}
}

func TestService_UseLog(t *testing.T) {
	s := NewService()
	s.UseLog(LogMW)
}

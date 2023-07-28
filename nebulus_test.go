package nebulus

import (
	"github.com/harley9293/nebulus/pkg/def"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go Run()
	time.Sleep(10 * time.Millisecond)
	m.Run()
	Shutdown()
}

type Service struct {
	def.DefaultHandler
}

func (m *Service) Print(req string) string {
	return "echo: " + req
}

func TestRegister(t *testing.T) {
	err := Register("echo", new(Service))
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	Destroy("echo")
	time.Sleep(10 * time.Millisecond)
}

func TestSend(t *testing.T) {
	err := Register("echo", new(Service))
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	time.Sleep(10 * time.Millisecond)
	Send("echo.Print", "hello world")

	Destroy("echo")
	time.Sleep(10 * time.Millisecond)
}

func TestCall(t *testing.T) {
	err := Register("echo", new(Service))
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	time.Sleep(10 * time.Millisecond)
	var rsp string
	err = Call("echo.Print", "hello world", &rsp)
	if err != nil {
		t.Fatal("Call() failed, err:" + err.Error())
	}
	if rsp != "echo: hello world" {
		t.Fatal("rsp != echo: hello world, rsp:" + rsp)
	}

	Destroy("echo")
	time.Sleep(10 * time.Millisecond)
}

func TestShutdown(t *testing.T) {
	Shutdown()
}

package nebulus

import (
	"github.com/harley9293/nebulus/pkg/def"
	"syscall"
	"testing"
	"time"
)

type Service struct {
	def.DefaultHandler
}

func (m *Service) Print(req string) string {
	return "echo: " + req
}

func TestRegister(t *testing.T) {
	defer Destroy("echo")
	err := Register("echo", new(Service))
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	go Run()
	time.Sleep(1 * time.Second)
}

func TestSend(t *testing.T) {
	defer Destroy("echo")
	err := Register("echo", new(Service))
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	go Run()
	time.Sleep(1 * time.Second)
	Send("echo.Print", "hello world")
}

func TestCall(t *testing.T) {
	defer Destroy("echo")
	err := Register("echo", new(Service))
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	go Run()
	time.Sleep(1 * time.Second)
	var rsp string
	err = Call("echo.Print", "hello world", &rsp)
	if err != nil {
		t.Fatal("Call() failed, err:" + err.Error())
	}
	if rsp != "echo: hello world" {
		t.Fatal("rsp != echo: hello world, rsp:" + rsp)
	}
}

func TestKill(t *testing.T) {
	svr.kill <- syscall.SIGTERM

	time.Sleep(1 * time.Second)
}

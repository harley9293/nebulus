package echo

import (
	"testing"

	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/internal/service"
)

func init() {
	nebulus.Register("Echo", &Service{})
	go nebulus.Run()
}

func TestService_Print(t *testing.T) {
	var rsp string
	err := service.Call("Echo.Print", "hello world", &rsp)
	if err != nil {
		t.Fatal(err)
	}

	if rsp != "echo: hello world" {
		t.Fatal(rsp)
	}
}

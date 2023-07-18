package httpd

import (
	"github.com/harley9293/nebulus/internal/service"
	"testing"
)

func TestNewHttpService(t *testing.T) {
	s := NewHttpService()
	if s == nil {
		t.Fatal("NewHttpService() failed")
	}

	type Req struct {
		Context string
	}

	type Rsp struct {
		Result string
	}

	err := s.AddHandler("GET", "/echo", func(req *Req, ctx *Context) Rsp {
		return Rsp{Result: req.Context}
	})
	if err != nil {
		t.Fatal("AddHandler() failed err:" + err.Error())
	}

	err = service.Register("http", s, "::8080")
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}
}

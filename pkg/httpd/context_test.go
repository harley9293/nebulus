package httpd

import (
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"testing"
)

func TestContext_CreateSession(t *testing.T) {
	ctx := &Context{service: NewHttpService(&Config{})}
	ctx.CreateSession("hello")
	if ctx.Session == nil {
		t.Fatal("create session error")
	}
}

func TestContext_Error(t *testing.T) {
	ctx := &Context{service: NewHttpService(&Config{})}
	ctx.Error(404, errors.New("test error"))
	if ctx.status != 404 || ctx.err.Error() != "test error" {
		t.Fatal("error error")
	}
}

func TestContext_Next_RspString(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30000", func(s *Service) {
		s.AddGlobalMiddleWare(LogMW, CookieMW, CorsMW)
		s.AddHandler("POST", "/test", func(req *EmptyTestStruct, ctx *Context) string {
			return "hello"
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30000")
	err := client.Post("/test", &EmptyTestStruct{})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

	if client.status != http.StatusOK {
		t.Fatal("doRequest() failed, status:" + string(rune(client.status)))
	}

	if client.strRsp != "hello" {
		t.Fatal("doRequest() failed, rsp:" + client.strRsp)
	}
}

func TestContext_Next_RspJsonMarshalError(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30001", func(s *Service) {
		type ComplexStructRsp struct {
			Ch chan int `json:"ch"`
		}
		s.AddGlobalMiddleWare(LogMW, CookieMW, CorsMW)
		s.AddHandler("POST", "/test", func(req *EmptyTestStruct, ctx *Context) ComplexStructRsp {
			return ComplexStructRsp{Ch: make(chan int)}
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30001")
	err := client.Post("/test", &EmptyTestStruct{})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

	if client.status != http.StatusInternalServerError {
		t.Fatal("doRequest() failed, status:" + string(rune(client.status)))
	}
}

package httpd

import (
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"testing"
)

func TestContext_CreateSession(t *testing.T) {
	ctx := &Context{service: NewHttpService()}
	ctx.CreateSession("hello")
	if ctx.Session == nil {
		t.Fatal("create session error")
	}
}

func TestContext_Error(t *testing.T) {
	ctx := &Context{service: NewHttpService()}
	ctx.Error(404, errors.New("test error"))
	if ctx.status != 404 || ctx.err.Error() != "test error" {
		t.Fatal("error error")
	}
}

func TestContext_Next(t *testing.T) {
	initTestEnv(t)

	echoRsp := &EchoRsp{}
	status, _, _ := doRequest(t, "POST", "/echo", "", &EchoReq{Content: "hello"}, echoRsp)
	if status == http.StatusOK {
		t.Fatal("status not ok")
	}

	req := &LoginReq{
		User: "harley9293",
		Pass: "123456",
	}
	rsp := &LoginRsp{}
	status, _, err := doRequest(t, "POST", "/login", "", req, rsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	status, _, err = doRequest(t, "POST", "/loginFail", "", req, rsp)
	if status == http.StatusOK {
		t.Fatal("status not ok")
	}
}

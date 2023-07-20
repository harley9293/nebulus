package httpd

import (
	"bytes"
	"encoding/json"
	"github.com/harley9293/nebulus"
	"io"
	"net/http"
	"testing"
)

type LoginReq struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginRsp struct {
	Token string `json:"token"`
}

func HandleLoginReq(req *LoginReq, ctx *Context) LoginRsp {
	rsp := LoginRsp{}
	ctx.CreateSession(req.User + req.Pass)

	if ctx.Session == nil {
		return rsp
	}

	rsp.Token = ctx.Session.id
	return rsp
}

type EchoReq struct {
	Content string `json:"content"`
}

type EchoRsp struct {
	Echo string `json:"echo"`
}

func HandleEchoReq(req *EchoReq, ctx *Context) EchoRsp {
	rsp := EchoRsp{}
	rsp.Echo = req.Content
	return rsp
}

func initTestEnv(t *testing.T) {
	s := NewHttpService()
	if s == nil {
		t.Fatal("NewHttpService() failed")
	}

	s.AddHandler("POST", "/login", HandleLoginReq, nil)
	s.AddHandler("POST", "/echo", HandleEchoReq, []MiddlewareFunc{SessionMiddleware})

	err := nebulus.Register("http", s, "127.0.0.1:36000")
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	go nebulus.Run()
}

func doRequest(t *testing.T, method, url string, req, rsp any) (status int, sessionID string) {
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal("json.Marshal() failed, err:" + err.Error())
	}
	data := bytes.NewBuffer(b)
	var resp *http.Response
	if method == "POST" {
		resp, err = http.Post("http://localhost:36000"+url, "application/json", data)
	} else {
		t.Fatal("method not support, method:" + method)
	}
	if err != nil {
		t.Fatal("http.Post() failed, err:" + err.Error())
	}
	status = resp.StatusCode
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal("Body.Close() failed, err:" + err.Error())
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil {
		t.Fatal("json.NewDecoder().Decode() failed, err:" + err.Error())
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_id" {
			sessionID = cookie.Value
			return
		}
	}
	return
}

func TestLogin(t *testing.T) {
	initTestEnv(t)

	req := &LoginReq{
		User: "harley9293",
		Pass: "123456",
	}
	rsp := &LoginRsp{}
	status, sessionID := doRequest(t, "POST", "/login", req, rsp)

	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	if rsp.Token == "" {
		t.Fatal("rsp.Token is empty")
	}

	if sessionID != rsp.Token {
		t.Fatal("sessionID != rsp.Token, sessionID:" + sessionID + ", rsp.Token:" + rsp.Token)
	}
}

func TestEchoWithLogin(t *testing.T) {
	initTestEnv(t)
}

func TestEchoWithoutLogin(t *testing.T) {
	initTestEnv(t)
}

func TestServiceFailed(t *testing.T) {
	initTestEnv(t)
}

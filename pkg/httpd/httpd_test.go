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

func doRequest(t *testing.T, method, url, session string, req, rsp any) (status int, sessionID string) {
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal("json.Marshal() failed, err:" + err.Error())
	}
	data := bytes.NewBuffer(b)
	var resp *http.Response
	request, err := http.NewRequest(method, "http://localhost:36000"+url, data)
	if err != nil {
		t.Fatal("http.NewRequest() failed, err:" + err.Error())
	}
	if session != "" {
		request.AddCookie(&http.Cookie{Name: "session_id", Value: session})
	}
	resp, err = http.DefaultClient.Do(request)
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

	if status == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(rsp)
		if err != nil {
			t.Fatal("json.NewDecoder().Decode() failed, err:" + err.Error())
		}
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
	status, sessionID := doRequest(t, "POST", "/login", "", req, rsp)

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
	req := &LoginReq{
		User: "harley9293",
		Pass: "123456",
	}
	rsp := &LoginRsp{}
	status, sessionID := doRequest(t, "POST", "/login", "", req, rsp)
	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	echoRsp := &EchoRsp{}
	status, sessionID = doRequest(t, "POST", "/echo", sessionID, &EchoReq{Content: "hello"}, echoRsp)
	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	if echoRsp.Echo != "hello" {
		t.Fatal("echoRsp.Echo != hello, echoRsp.Echo:" + echoRsp.Echo)
	}
}

func TestEchoWithoutLogin(t *testing.T) {
	initTestEnv(t)
	echoRsp := &EchoRsp{}
	status, _ := doRequest(t, "POST", "/echo", "", &EchoReq{Content: "hello"}, echoRsp)
	if status != http.StatusUnauthorized {
		t.Fatal("status not 401, status:" + string(rune(status)))
	}
}

func TestServiceFailed(t *testing.T) {
	initTestEnv(t)
}

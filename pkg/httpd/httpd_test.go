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

var serviceStart = false

type PanicReq struct {
}

type PanicRsp struct {
}

func HandlePanicReq(req *PanicReq, ctx *Context) PanicRsp {
	panic("test panic")
}

func initTestEnv(t *testing.T) {
	if serviceStart {
		return
	}

	s := DefaultHttpService()
	if s == nil {
		t.Fatal("NewHttpService() failed")
	}

	s.AddHandler("POST", "/login", HandleLoginReq)
	s.AddHandler("POST", "/echo", HandleEchoReq, AuthMW)
	s.AddHandler("POST", "/panic", HandlePanicReq)

	err := nebulus.Register("http", s, "127.0.0.1:36000")
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	serviceStart = true
	go nebulus.Run()
}

func doRequest(t *testing.T, method, url, session string, req, rsp any) (status int, sessionID string, err error) {
	b, err := json.Marshal(req)
	if err != nil {
		return 0, "", err
	}
	data := bytes.NewBuffer(b)
	var resp *http.Response
	request, err := http.NewRequest(method, "http://localhost:36000"+url, data)
	if err != nil {
		return 0, "", err
	}
	if session != "" {
		request.AddCookie(&http.Cookie{Name: "token", Value: session})
	}
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return 0, "", err
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
			return 0, "", err
		}
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
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
	status, sessionID, err := doRequest(t, "POST", "/login", "", req, rsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

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
	status, sessionID, err := doRequest(t, "POST", "/login", "", req, rsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	echoRsp := &EchoRsp{}
	status, sessionID, err = doRequest(t, "POST", "/echo", sessionID, &EchoReq{Content: "hello"}, echoRsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

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
	status, _, err := doRequest(t, "POST", "/echo", "", &EchoReq{Content: "hello"}, echoRsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

	if status != http.StatusUnauthorized {
		t.Fatal("status not 401, status:" + string(rune(status)))
	}
}

func TestServiceFailed(t *testing.T) {
	initTestEnv(t)

	req := &LoginReq{
		User: "harley9293",
		Pass: "123456",
	}
	rsp := &LoginRsp{}
	status, sessionID, err := doRequest(t, "POST", "/login", "", req, rsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	panicRsp := &PanicRsp{}
	_, _, err = doRequest(t, "POST", "/panic", "", &PanicReq{}, panicRsp)
	if err == nil {
		t.Fatal("doRequest() failed, err is nil")
	}

	echoRsp := &EchoRsp{}
	status, sessionID, err = doRequest(t, "POST", "/echo", sessionID, &EchoReq{Content: "hello"}, echoRsp)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}

	if status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(status)))
	}

	if echoRsp.Echo != "hello" {
		t.Fatal("echoRsp.Echo != hello, echoRsp.Echo:" + echoRsp.Echo)
	}
}

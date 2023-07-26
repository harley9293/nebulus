package httpd

import (
	"bytes"
	"encoding/json"
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/errors"
	"io"
	"net/http"
	"testing"
	"time"
)

type LoginReq struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginErrReq struct {
	User string `json:"user"`
	Pass int    `json:"pass"`
}

type LoginRsp struct {
	Token string `json:"token"`
}

type PackRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func HandleLoginReq(req *LoginReq, ctx *Context) LoginRsp {
	rsp := LoginRsp{}
	if req.User == "111111" {
		ctx.Error(http.StatusInternalServerError, errors.New("test error"))
		return rsp
	}

	ctx.CreateSession(req.User + req.Pass)

	if ctx.Session == nil {
		return rsp
	}

	rsp.Token = ctx.Session.id
	return rsp
}

func HandleLoginReqFail(req *LoginReq, ctx *Context) func() {
	return func() {}
}

func HandleMockResponse(req *LoginReq, ctx *Context) LoginRsp {
	rsp := LoginRsp{}
	ctx.w = &mockResponseWriter{}
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

type PanicReq struct {
}

type PanicRsp struct {
}

func HandlePanicReq(req *PanicReq, ctx *Context) PanicRsp {
	panic("test panic")
}

var serviceTest *Service = nil

func initTestEnv(t *testing.T) {
	if serviceTest != nil {
		return
	}

	s := NewHttpService()
	if s == nil {
		t.Fatal("NewHttpService() failed")
	}

	s.AddGlobalMiddleWare(LogMW, CookieMW, CorsMW)
	s.AddHandler("POST", "/login", HandleLoginReq)
	s.AddHandler("POST", "/loginFail", HandleLoginReqFail)
	s.AddHandler("POST", "/echo", HandleEchoReq, AuthMW)
	s.AddHandler("POST", "/panic", HandlePanicReq)
	s.AddHandler("POST", "/mockResponse", HandleMockResponse)
	s.AddHandler("POST", "/loginPack", HandleLoginReq, RspPackMW)

	err := nebulus.Register("http", s, "127.0.0.1:36000")
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	serviceTest = s
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
	if resp != nil {
		status = resp.StatusCode
	}
	if err != nil {
		return 0, "", err
	}
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

func TestService_AddHandler(t *testing.T) {
	defer func() { recover() }()

	service := NewHttpService()
	service.AddHandler("GET", "/test", FailFuncTest1)
}

func TestService_OnInit(t *testing.T) {
	service := NewHttpService()
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

func TestServiceFailed2(t *testing.T) {
	initTestEnv(t)

	err := serviceTest.srv.Close()
	if err != nil {
		t.Fatal("srv.Close() failed, err:" + err.Error())
	}

	time.Sleep(1 * time.Second)

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
}

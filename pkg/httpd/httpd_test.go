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

func TestInitEnv(t *testing.T) {
	s := NewHttpService()
	if s == nil {
		t.Fatal("NewHttpService() failed")
	}

	err := s.AddHandler("POST", "/login", HandleLoginReq, nil)
	if err != nil {
		t.Fatal("AddHandler() /login failed err:" + err.Error())
	}

	err = s.AddHandler("GET", "/echo", HandleEchoReq, []MiddlewareFunc{SessionMiddleware})
	if err != nil {
		t.Fatal("AddHandler() /echo failed err:" + err.Error())
	}

	err = nebulus.Register("http", s, "127.0.0.1:36000")
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}

	go nebulus.Run()
}

func TestNewHttpService(t *testing.T) {
	TestInitEnv(t)

	req := LoginReq{
		User: "harley9293",
		Pass: "123456",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal("json.Marshal() failed, err:" + err.Error())
	}
	data := bytes.NewBuffer(b)
	resp, err := http.Post("http://localhost:36000/login", "application/json", data)
	if err != nil {
		t.Fatal("http.Post() failed, err:" + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal("Body.Close() failed, err:" + err.Error())
		}
	}(resp.Body)

	rsp := LoginRsp{}
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil {
		t.Fatal("json.NewDecoder().Decode() failed, err:" + err.Error())
	}

	if rsp.Token == "" {
		t.Fatal("rsp.Token is empty")
	}
}

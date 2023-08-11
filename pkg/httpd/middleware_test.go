package httpd

import (
	"bytes"
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"strconv"
	"testing"
)

type mockResponseWriter struct {
	header http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	// 可以在这里返回任何你需要的错误。
	return 0, errors.New("mock error")
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
}

func TestRouterMW_NotFound(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30006", func(s *Service) {
		s.AddHandler("POST", "/test", func(testStruct *EmptyTestStruct, ctx *Context) string {
			return ""
		})
	})
	defer nebulus.Destroy("Test")
	client := NewClient("http://127.0.0.1:30006")
	err := client.Post("/test2", &EmptyTestStruct{})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusNotFound {
		t.Fatal("status not statusNotFound, status:" + string(rune(client.status)))
	}
}

func TestRouterMW_RequestDecodeError(t *testing.T) {
	type TestReq struct {
		Test int `json:"test"`
	}
	type ErrorReq struct {
		Test string `json:"test"`
	}
	NewTestHttpService("Test", "127.0.0.1:30008", func(s *Service) {
		s.AddHandler("POST", "/test", func(req *TestReq, ctx *Context) string {
			return ""
		})
		s.AddHandler("GET", "/test2", func(req *TestReq, ctx *Context) string {
			return ""
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30008")
	err := client.Post("/test", &ErrorReq{Test: "1"})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusBadRequest {
		t.Fatal("status not statusBadRequest, status:" + string(rune(client.status)))
	}

	err = client.Get("/test2", map[string]string{"test": "1"})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusBadRequest {
		t.Fatal("status not statusBadRequest, status:" + string(rune(client.status)))
	}
}

func TestRouterMW_UrlParamDecodeSuccess(t *testing.T) {
	type TestReq struct {
		Test string `json:"test"`
	}
	NewTestHttpService("Test", "127.0.0.1:30009", func(s *Service) {
		s.AddHandler("GET", "/test", func(req *TestReq, ctx *Context) string {
			return "hello " + req.Test
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30009")
	err := client.Get("/test", map[string]string{"test": "world"})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not statusOK, status:" + string(rune(client.status)))
	}
	if client.strRsp != "hello world" {
		t.Fatal("rsp not ok, rsp:" + client.strRsp)
	}
}

func TestRouterMW_UrlParamDecodeError(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30010", func(s *Service) {
		s.AddHandler("GET", "/test", func(req *EmptyTestStruct, ctx *Context) string {
			return "hello"
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30010")
	client.method = "GET"
	client.url = client.host + "/test" + "?" + "Hello%2Gworld"
	client.body = bytes.NewBuffer([]byte{})
	err := client.do()
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusBadRequest {
		t.Fatal("status not statusBadRequest, status:" + string(rune(client.status)))
	}
}

func TestResponseMW_WriteResponseError(t *testing.T) {
	type TestRsp struct {
		Test string `json:"test"`
	}
	NewTestHttpService("Test", "127.0.0.1:30005", func(s *Service) {
		s.AddHandler("POST", "/test", func(testStruct *EmptyTestStruct, ctx *Context) string {
			ctx.w = &mockResponseWriter{header: http.Header{}}
			return ""
		})
		s.AddHandler("POST", "/test2", func(ctx *Context) (rsp TestRsp) {
			ctx.w = &mockResponseWriter{header: http.Header{}}
			ctx.w.Header().Set("Content-Type", "application/json")
			return
		})
	})
	defer nebulus.Destroy("Test")
	client := NewClient("http://127.0.0.1:30005")
	_ = client.Post("/test", &EmptyTestStruct{})
	_ = client.Post("/test2", nil)
}

func TestAuthMW_Fail(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30003", func(s *Service) {
		s.AddGlobalMiddleWare(CookieMW, CorsMW)
		s.AddHandler("POST", "/test", func(testStruct *EmptyTestStruct, ctx *Context) string {
			return "test"
		}, AuthMW)
	})
	defer nebulus.Destroy("Test")
	client := NewClient("http://127.0.0.1:30003")
	err := client.Post("/test", &EmptyTestStruct{})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusUnauthorized {
		t.Fatal("status not statusUnauthorized, status:" + string(rune(client.status)))
	}
}

func TestAuthMW_Success(t *testing.T) {
	type LoginReq struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}
	NewTestHttpService("Test", "127.0.0.1:30002", func(s *Service) {
		s.AddGlobalMiddleWare(CookieMW, CorsMW)
		s.AddHandler("POST", "/login", func(req *LoginReq, ctx *Context) string {
			ctx.CreateSession(req.User + req.Pass)
			return ""
		})
		s.AddHandler("POST", "/test", func(testStruct *EmptyTestStruct, ctx *Context) string {
			return "test"
		}, AuthMW)
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30002")
	err := client.Post("/login", &LoginReq{"harley9293", "123456"})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}

	err = client.Post("/test", &EmptyTestStruct{})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}
	if client.strRsp != "test" {
		t.Fatal("rsp not ok, rsp:" + client.strRsp)
	}
}

func TestRspPackMW_Success(t *testing.T) {
	type TestRsp struct {
		Test string `json:"test"`
	}

	type TestPackRsp struct {
		Code int     `json:"code"`
		Msg  string  `json:"msg"`
		Data TestRsp `json:"data"`
	}
	NewTestHttpService("Test", "127.0.0.1:31001", func(s *Service) {
		s.AddGlobalMiddleWare(RspPackMW)
		s.AddHandler("POST", "/test", func(ctx *Context) (rsp TestRsp) {
			rsp.Test = "hello world!"
			return
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:31001")
	err := client.Post("/test", nil)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}

	var rsp TestPackRsp
	err = client.jsonRsp.Decode(&rsp)
	if err != nil {
		t.Fatal("decode rsp failed, err:" + err.Error())
	}
	if rsp.Code != http.StatusOK {
		t.Fatal("rsp.Code not ok, rsp.Code:" + strconv.Itoa(rsp.Code))
	}
	if rsp.Msg != "success" {
		t.Fatal("rsp.Msg not ok, rsp.Msg:" + rsp.Msg)
	}
	if rsp.Data.Test != "hello world!" {
		t.Fatal("rsp.Data.Test not ok, rsp.Data.Test:" + rsp.Data.Test)
	}
}

func TestRspPackMW_Fail(t *testing.T) {
	type TestPackRsp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	NewTestHttpService("Test", "127.0.0.1:31002", func(s *Service) {
		s.AddGlobalMiddleWare(RspPackMW)
		s.AddHandler("POST", "/test", func(ctx *Context) {
			ctx.Error(http.StatusInternalServerError, errors.New("test error"))
			return
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:31002")
	err := client.Post("/test", nil)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}

	var rsp TestPackRsp
	err = client.jsonRsp.Decode(&rsp)
	if err != nil {
		t.Fatal("decode rsp failed, err:" + err.Error())
	}
	if rsp.Code != http.StatusInternalServerError {
		t.Fatal("rsp.Code not ok, rsp.Code:" + strconv.Itoa(rsp.Code))
	}
	if rsp.Msg != "test error" {
		t.Fatal("rsp.Msg not ok, rsp.Msg:" + rsp.Msg)
	}
}

package httpd

import (
	"bytes"
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"testing"
)

type mockResponseWriter struct {
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
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

func TestRouterMW_MethodNotAllowed(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30007", func(s *Service) {
		s.AddHandler("POST", "/test", func(testStruct *EmptyTestStruct, ctx *Context) string {
			return ""
		})
	})
	defer nebulus.Destroy("Test")
	client := NewClient("http://127.0.0.1:30007")
	err := client.Get("/test", nil)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusMethodNotAllowed {
		t.Fatal("status not statusMethodNotAllowed, status:" + string(rune(client.status)))
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
	NewTestHttpService("Test", "127.0.0.1:30005", func(s *Service) {
		s.AddHandler("POST", "/test", func(testStruct *EmptyTestStruct, ctx *Context) string {
			ctx.w = &mockResponseWriter{}
			return ""
		})
	})
	defer nebulus.Destroy("Test")
	client := NewClient("http://127.0.0.1:30005")
	_ = client.Post("/test", &EmptyTestStruct{})
}

func TestAuthMW_Fail(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:30003", func(s *Service) {
		s.AddGlobalMiddleWare(LogMW, CookieMW, CorsMW)
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
		s.AddGlobalMiddleWare(LogMW, CookieMW, CorsMW)
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

func TestRspPackMW(t *testing.T) {
	type TestReq struct {
		Ok bool `json:"ok"`
	}
	NewTestHttpService("Test", "127.0.0.1:30004", func(s *Service) {
		s.AddHandler("POST", "/test", func(req *TestReq, ctx *Context) string {
			if req.Ok {
				return "true"
			}
			ctx.Error(http.StatusInternalServerError, errors.New("test error"))
			return ""
		}, RspPackMW)
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:30004")
	err := client.Post("/test", &TestReq{Ok: true})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusOK {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}
	result := make(map[string]any)
	err = client.jsonRsp.Decode(&result)
	if err != nil {
		t.Fatal("decode rsp failed, err:" + err.Error())
	}
	if result["code"].(float64) != 200 || result["msg"].(string) != "success" || result["data"].(string) != "true" {
		t.Fatal("rsp not ok")
	}

	err = client.Post("/test", &TestReq{Ok: false})
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.status != http.StatusInternalServerError {
		t.Fatal("status not ok, status:" + string(rune(client.status)))
	}
}

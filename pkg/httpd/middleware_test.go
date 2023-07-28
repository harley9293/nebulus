package httpd

import (
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

func TestRouterMW(t *testing.T) {
	//initTestEnv(t)
	//
	//req := &LoginReq{
	//	User: "harley9293",
	//	Pass: "123456",
	//}
	//rsp := &LoginRsp{}
	//
	//status, _, err := doRequest(t, "GET", "/base", "", req, rsp)
	//if status != http.StatusNotFound {
	//	t.Fatalf("expect status %d, got %d, err: %s", http.StatusNotFound, status, err.Error())
	//}
	//
	//status, _, _ = doRequest(t, "GET", "/login", "", req, rsp)
	//if status != http.StatusMethodNotAllowed {
	//	t.Fatalf("expect status %d, got %d", http.StatusMethodNotAllowed, status)
	//}
	//
	//status, _, _ = doRequest(t, "POST", "/login", "", &LoginErrReq{"harley9293", 123456}, rsp)
	//if status != http.StatusBadRequest {
	//	t.Fatalf("expect status %d, got %d", http.StatusBadRequest, status)
	//}
	//
	//rsp = &LoginRsp{}
	//status, err = doRequestGet("/loginGet", map[string]string{"user": "admin", "pass": "123456"}, rsp, "")
	//if err != nil {
	//	t.Fatalf("doRequestGet() failed, err:" + err.Error())
	//}
	//
	//if status != http.StatusOK {
	//	t.Fatalf("expect status %d, got %d", http.StatusOK, status)
	//}
	//
	//if rsp.Token == "" {
	//	t.Fatalf("expect token not empty, got empty")
	//}
	//
	//status, _ = doRequestGet("/loginGet", map[string]string{"user": "admin", "pass": "123456"}, rsp, "Hello%2Gworld")
	//if status != http.StatusBadRequest {
	//	t.Fatalf("expect status %d, got %d", http.StatusBadRequest, status)
	//}
	//
	//status, _ = doRequestGet("/loginGetFail", map[string]string{"user": "admin", "pass": "123456"}, rsp, "")
	//if status != http.StatusBadRequest {
	//	t.Fatalf("expect status %d, got %d", http.StatusBadRequest, status)
	//}
}

func TestResponseMW(t *testing.T) {
	//initTestEnv(t)
	//
	//req := &LoginReq{
	//	User: "harley9293",
	//	Pass: "123456",
	//}
	//rsp := &LoginRsp{}
	//_, _, err := doRequest(t, "POST", "/mockResponse", "", req, rsp)
	//if err == nil {
	//	t.Fatal("expect error, got nil")
	//}
}

func TestAuthMW(t *testing.T) {
	//initTestEnv(t)
	//req := &LoginReq{
	//	User: "harley9293",
	//	Pass: "123456",
	//}
	//rsp := &LoginRsp{}
	//status, sessionID, err := doRequest(t, "POST", "/login", "", req, rsp)
	//if err != nil {
	//	t.Fatal("doRequest() failed, err:" + err.Error())
	//}
	//
	//if status != http.StatusOK {
	//	t.Fatal("status not ok, status:" + string(rune(status)))
	//}
	//
	//echoRsp := &EchoRsp{}
	//status, sessionID, err = doRequest(t, "POST", "/echo", sessionID, &EchoReq{Content: "hello"}, echoRsp)
	//if err != nil {
	//	t.Fatal("doRequest() failed, err:" + err.Error())
	//}
	//
	//if status != http.StatusOK {
	//	t.Fatal("status not ok, status:" + string(rune(status)))
	//}
	//
	//if echoRsp.Echo != "hello" {
	//	t.Fatal("echoRsp.Echo != hello, echoRsp.Echo:" + echoRsp.Echo)
	//}
}

func TestRspPackMW(t *testing.T) {
	//initTestEnv(t)
	//req := &LoginReq{
	//	User: "harley9293",
	//	Pass: "123456",
	//}
	//rsp := &PackRsp{}
	//status, _, err := doRequest(t, "POST", "/loginPack", "", req, rsp)
	//if status != 200 {
	//	t.Fatalf("expect status %d, got %d, err: %s", 200, status, err.Error())
	//}
	//if err != nil {
	//	t.Fatalf("expect err nil, got %s", err.Error())
	//}
	//
	//status, _, err = doRequest(t, "POST", "/loginPack", "", &LoginReq{"111111", "123456"}, rsp)
	//if status != http.StatusInternalServerError {
	//	t.Fatalf("expect status %d, got %d, err: %s", http.StatusInternalServerError, status, err.Error())
	//}
}

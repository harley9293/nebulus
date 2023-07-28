package httpd

import (
	"bytes"
	"github.com/harley9293/nebulus"
	"testing"
)

func TestClient_Post_MarshalError(t *testing.T) {
	NewTestHttpService("Test", "127.0.0.1:36013", func(s *Service) {
		s.AddHandler("POST", "/test", func(req *EmptyTestStruct, ctx *Context) string {
			return "hello world"
		})
	})
	defer nebulus.Destroy("Test")

	client := NewClient("http://127.0.0.1:36013")
	err := client.Post("/test", make(chan int))
	if err == nil {
		t.Fatal("doRequest() failed, err is nil")
	}

	err = client.Post("/test", nil)
	if err != nil {
		t.Fatal("doRequest() failed, err:" + err.Error())
	}
	if client.strRsp != "hello world" {
		t.Fatal("doRequest() failed, rsp:" + client.strRsp)
	}
}

func TestClient_Do_NewRequestError(t *testing.T) {
	client := NewClient("http://127.0.0.1")
	client.method = "TEST"
	client.body = bytes.NewBuffer([]byte{})
	err := client.do()
	if err == nil {
		t.Fatal("do() failed, err is nil")
	}
}

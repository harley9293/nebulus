package httpd

import (
	"github.com/harley9293/nebulus/internal/service"
	"net/http"
	"testing"
)

func TestNewHttpService(t *testing.T) {
	s := NewHttpService()
	if s == nil {
		t.Fatal("NewHttpService() failed")
	}

	s.AddHandler("/echo", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("echo"))
		if err != nil {
			t.Fatal("http write failed")
		}
	})

	err := service.Register("http", s, "::8080")
	if err != nil {
		t.Fatal("Register() failed, err:" + err.Error())
	}
}

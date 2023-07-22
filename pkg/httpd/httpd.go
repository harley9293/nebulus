package httpd

import (
	"github.com/harley9293/nebulus/pkg/def"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"reflect"
	"time"
)

var InitArgsSizeError = errors.New("http init args size error, got:%d")
var InitArgsTypeError = errors.New("http init args type error, got:%T")

type MiddlewareFunc func(*Context)

type Service struct {
	def.DefaultHandler

	srv *http.Server
	err chan error

	hm *handlerMng
}

func NewHttpService() *Service {
	return &Service{hm: newHandlerMng()}
}

func (m *Service) AddHandler(method, path string, f any, middleware []MiddlewareFunc) {
	err := m.hm.add(method, path, f, middleware)
	if err != nil {
		panic(err)
	}
}

func (m *Service) AddGlobalMiddleWare(f MiddlewareFunc) {
	m.hm.gmw = append(m.hm.gmw, f)
}

func (m *Service) OnInit(args ...any) error {
	if len(args) != 1 {
		return InitArgsSizeError.Fill(len(args))
	}

	address, ok := args[0].(string)
	if !ok {
		return InitArgsTypeError.Fill(reflect.TypeOf(args[0]))
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", m.hm.handler)
	m.srv = &http.Server{Addr: address, Handler: serveMux}
	m.err = make(chan error, 1)

	go func() {
		m.err <- m.srv.ListenAndServe()
	}()
	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-m.err:
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

func (m *Service) OnTick() {
	select {
	case err, ok := <-m.err:
		if ok {
			close(m.err)
			panic(err)
		}
	}
}

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

	globalMiddlewares []MiddlewareFunc
	hm                *handlerMng
	sm                *sessionMng
}

func DefaultHttpService() *Service {
	service := NewHttpService()
	service.AddGlobalMiddleWare(PreRequestMW, PreResponseMW)
	return service
}

func NewHttpService() *Service {
	return &Service{hm: newHandlerMng(), sm: newSessionMng()}
}

func (m *Service) AddHandler(method, path string, f any, middleware ...MiddlewareFunc) {
	err := m.hm.add(method, path, f, middleware)
	if err != nil {
		panic(err)
	}
}

func (m *Service) AddGlobalMiddleWare(f ...MiddlewareFunc) {
	m.globalMiddlewares = f
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
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			r:           r,
			w:           w,
			service:     m,
			status:      http.StatusOK,
			middlewares: m.globalMiddlewares,
		}
		c.Next()
	})
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

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

	srv      *http.Server
	address  string
	serveMux *http.ServeMux
	err      chan error

	globalMiddlewares []MiddlewareFunc
	router            *router
	sm                *sessionMng
	cfg               *Config
}

func NewService(config *Config) *Service {
	config.Fill()
	service := &Service{router: newRouter(), sm: newSessionMng(config.SType, config.SExpireTime, config.Redis), cfg: config}
	service.AddGlobalMiddleWare(responseMW, routerMW)
	return service
}

func (m *Service) AddHandler(method, path string, f any, middleware ...MiddlewareFunc) {
	err := m.router.add(method, path, f, middleware...)
	if err != nil {
		panic(err)
	}
}

func (m *Service) AddGlobalMiddleWare(f ...MiddlewareFunc) {
	m.globalMiddlewares = append(m.globalMiddlewares, f...)
}

func (m *Service) OnInit(args ...any) error {
	if len(args) != 1 {
		return InitArgsSizeError.Fill(len(args))
	}

	address, ok := args[0].(string)
	if !ok {
		return InitArgsTypeError.Fill(reflect.TypeOf(args[0]))
	}
	m.address = address

	m.serveMux = http.NewServeMux()
	m.serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			r:           r,
			w:           w,
			service:     m,
			status:      http.StatusOK,
			middlewares: m.globalMiddlewares,
		}
		c.Next()
	})
	m.srv = &http.Server{Addr: m.address, Handler: m.serveMux}
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
	case err := <-m.err:
		panic(err)
	default:
	}
}

func (m *Service) OnPanic() {
	m.srv = &http.Server{Addr: m.address, Handler: m.serveMux}
	go func() {
		m.err <- m.srv.ListenAndServe()
	}()
}

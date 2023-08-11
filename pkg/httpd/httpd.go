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

	logMiddleware     MiddlewareFunc
	globalMiddlewares []MiddlewareFunc
	router            *router

	sessionMap  map[string]def.Session
	baseSession def.Session
}

func NewService() *Service {
	service := &Service{
		logMiddleware: LogMW,
		router:        newRouter(),
		sessionMap:    make(map[string]def.Session),
		baseSession:   &defaultSession{cfgExpireTime: 24 * time.Hour},
	}
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

func (m *Service) UseSession(session def.Session) {
	m.baseSession = session
}

func (m *Service) UseLog(f MiddlewareFunc) {
	m.logMiddleware = f
}

func (m *Service) GetSession(id string) def.Session {
	if session, ok := m.sessionMap[id]; ok {
		if session.IsExpired() {
			delete(m.sessionMap, id)
			return nil
		} else {
			return session
		}
	} else {
		return nil
	}
}

func (m *Service) NewSession(id string) def.Session {
	session := m.baseSession.New(id)
	m.sessionMap[session.ID()] = session
	return session
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
			middlewares: append([]MiddlewareFunc{m.logMiddleware}, m.globalMiddlewares...),
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

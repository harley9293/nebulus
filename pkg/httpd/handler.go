package httpd

import (
	"encoding/json"
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"reflect"
)

var HandlerTypeError = errors.New("add handler type error, got:%T")
var HandlerParamSizeError = errors.New("handler num in: %d, num out: %d")
var HandlerParamPointerError = errors.New("handler param must be pointer")
var HandlerSecondParamTypeError = errors.New("handler second param must be *Context")

type handlerData struct {
	path        string
	method      string
	handler     reflect.Value
	middlewares []MiddlewareFunc
}

type handlerMng struct {
	data map[string]*handlerData

	sm *sessionMng
}

func newHandlerMng() *handlerMng {
	return &handlerMng{
		data: make(map[string]*handlerData),
		sm:   newSessionMng(),
	}
}

func handlerVerify(value reflect.Value) error {
	if value.Type().Kind() != reflect.Func {
		return HandlerTypeError.Fill(reflect.TypeOf(value))
	}

	if value.Type().NumIn() != 2 || value.Type().NumOut() != 1 {
		return HandlerParamSizeError.Fill(value.Type().NumIn(), value.Type().NumOut())
	}

	if value.Type().In(0).Kind() != reflect.Ptr || value.Type().In(1).Kind() != reflect.Ptr {
		return HandlerParamPointerError
	}

	if value.Type().In(1) != reflect.TypeOf(&Context{}) {
		return HandlerSecondParamTypeError
	}

	return nil
}

func (m *handlerMng) add(method, path string, f any, middlewares []MiddlewareFunc) error {
	fn := reflect.ValueOf(f)
	err := handlerVerify(fn)
	if err != nil {
		return err
	}

	m.data[path] = &handlerData{
		path:        path,
		method:      method,
		handler:     fn,
		middlewares: middlewares,
	}

	return nil
}

func (m *handlerMng) handler(w http.ResponseWriter, r *http.Request) {
	h, ok := m.data[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.Method != h.method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	arg := reflect.New(h.handler.Type().In(0))
	err := json.NewDecoder(r.Body).Decode(arg.Interface())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Debug("url: %s, req: %+v", r.URL, arg.Elem().Interface())

	c := &Context{
		Session: nil,
		r:       r,
		w:       w,
		sm:      m.sm,
		h:       h,
	}

	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		c.Session = m.sm.get(sessionCookie.Value)
	}
	c.values = []reflect.Value{arg.Elem(), reflect.ValueOf(c)}

	err = c.Next()
	if err != nil {
		http.Error(c.w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("url: %s, err: %s", r.URL, err.Error())
		return
	}

	if c.Session != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: c.Session.id,
		})
	}
}

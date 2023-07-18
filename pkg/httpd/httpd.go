package httpd

import (
	"encoding/json"
	"github.com/harley9293/nebulus/pkg/def"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"reflect"
)

var InitArgsSizeError = errors.New("http init args size error, got:%d")
var InitArgsTypeError = errors.New("http init args type error, got:%T")
var HandlerTypeError = errors.New("add handler type error, got:%T")
var HandlerParamSizeError = errors.New("handler num in: %d, num out: %d")

type handlerData struct {
	path    string
	method  string
	handler reflect.Value
	argType reflect.Type
}

type Service struct {
	def.DefaultHandler

	srv        *http.Server
	handlerMap map[string]*handlerData
	err        chan error
}

func NewHttpService() *Service {
	return &Service{handlerMap: map[string]*handlerData{}}
}

func (m *Service) AddHandler(method, path string, f any) error {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return HandlerTypeError.Fill(reflect.TypeOf(f))
	}

	fn := reflect.ValueOf(f)
	if fn.Type().NumIn() != 1 || fn.Type().NumOut() != 1 {
		return HandlerParamSizeError.Fill(fn.Type().NumIn(), fn.Type().NumOut())
	}
	argType := fn.Type().In(0)

	m.handlerMap[path] = &handlerData{
		path:    path,
		method:  method,
		handler: fn,
		argType: argType,
	}

	return nil
}

func (m *Service) handler(w http.ResponseWriter, r *http.Request) {
	h, ok := m.handlerMap[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.Method != h.method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	arg := reflect.New(h.argType)

	err := json.NewDecoder(r.Body).Decode(arg.Interface())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	result := h.handler.Call([]reflect.Value{arg.Elem()})

	err = json.NewEncoder(w).Encode(result[0].Interface())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
	serveMux.HandleFunc("/", m.handler)
	m.srv = &http.Server{Addr: address, Handler: serveMux}

	m.err = make(chan error)

	go func() {
		m.err <- m.srv.ListenAndServe()
	}()

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

package httpd

import (
	"encoding/json"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"reflect"
)

var HandlerTypeError = errors.New("add handler type error, got:%T")
var HandlerParamSizeError = errors.New("handler num in: %d, num out: %d")

type handlerData struct {
	path    string
	method  string
	handler reflect.Value
}

type handlerMng struct {
	data map[string]*handlerData
}

func handlerVerify(value reflect.Value) error {
	if value.Type().Kind() != reflect.Func {
		return HandlerTypeError.Fill(reflect.TypeOf(value))
	}

	if value.Type().NumIn() != 1 || value.Type().NumOut() != 1 {
		return HandlerParamSizeError.Fill(value.Type().NumIn(), value.Type().NumOut())
	}

	return nil
}

func (m *handlerMng) add(method, path string, f any) error {
	fn := reflect.ValueOf(f)
	err := handlerVerify(fn)
	if err != nil {
		return err
	}

	m.data[path] = &handlerData{
		path:    path,
		method:  method,
		handler: fn,
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

	result := h.handler.Call([]reflect.Value{arg.Elem()})

	err = json.NewEncoder(w).Encode(result[0].Interface())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

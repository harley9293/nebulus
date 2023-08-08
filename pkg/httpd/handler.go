package httpd

import (
	"github.com/harley9293/nebulus/pkg/errors"
	"reflect"
)

var HandlerTypeError = errors.New("add handler type error, got:%T")
var HandlerParamPointerError = errors.New("handler param must be pointer")
var HandlerSecondParamTypeError = errors.New("handler second param must be *Context")
var HandlerRepeatedError = errors.New("handler path already exists, path:%s")

type handlerData struct {
	path        string
	method      string
	handler     reflect.Value
	middlewares []MiddlewareFunc
}

type handlerMng struct {
	data map[string]*handlerData
}

func newHandlerMng() *handlerMng {
	return &handlerMng{
		data: make(map[string]*handlerData),
	}
}

func handlerVerify(value reflect.Value) error {
	if value.Type().Kind() != reflect.Func {
		return HandlerTypeError.Fill(reflect.TypeOf(value))
	}

	for i := 0; i < value.Type().NumIn(); i++ {
		if value.Type().In(i).Kind() != reflect.Ptr {
			return HandlerParamPointerError
		}
	}

	if value.Type().In(value.Type().NumIn()-1) != reflect.TypeOf(&Context{}) {
		return HandlerSecondParamTypeError
	}

	return nil
}

func (m *handlerMng) add(method, path string, f any, middlewares ...MiddlewareFunc) error {
	if _, ok := m.data[path]; ok {
		return HandlerRepeatedError.Fill(path)
	}

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

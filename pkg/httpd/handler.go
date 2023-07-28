package httpd

import (
	"github.com/harley9293/nebulus/pkg/errors"
	"reflect"
)

var HandlerTypeError = errors.New("add handler type error, got:%T")
var HandlerParamSizeError = errors.New("handler num in: %d, num out: %d")
var HandlerParamPointerError = errors.New("handler param must be pointer")
var HandlerSecondParamTypeError = errors.New("handler second param must be *Context")
var HandlerRepeatedError = errors.New("handler path already exists, path:%s")
var HandlerReturnTypeError = errors.New("handler return type error, got:%T")

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

	if value.Type().NumIn() != 2 || value.Type().NumOut() != 1 {
		return HandlerParamSizeError.Fill(value.Type().NumIn(), value.Type().NumOut())
	}

	if value.Type().In(0).Kind() != reflect.Ptr || value.Type().In(1).Kind() != reflect.Ptr {
		return HandlerParamPointerError
	}

	if value.Type().In(1) != reflect.TypeOf(&Context{}) {
		return HandlerSecondParamTypeError
	}

	if value.Type().Out(0).Kind() != reflect.String && value.Type().Out(0).Kind() != reflect.Struct {
		return HandlerReturnTypeError.Fill(value.Type().Out(0))
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

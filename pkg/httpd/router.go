package httpd

import (
	"github.com/harley9293/nebulus/pkg/errors"
	"reflect"
)

var HandlerTypeError = errors.New("add handler type error, got:%T")
var HandlerParamPointerError = errors.New("handler param must be pointer")
var HandlerSecondParamTypeError = errors.New("handler second param must be *Context")
var HandlerRepeatedError = errors.New("handler path already exists,method:%s, path:%s")

type routes struct {
	path        string
	method      string
	fn          reflect.Value
	middlewares []MiddlewareFunc
}

type router struct {
	data []*routes
}

func newRouter() *router {
	return &router{}
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

func (m *router) add(method, path string, f any, middlewares ...MiddlewareFunc) error {
	for _, v := range m.data {
		if v.method == method && v.path == path {
			return HandlerRepeatedError.Fill(method, path)
		}
	}

	fn := reflect.ValueOf(f)
	err := handlerVerify(fn)
	if err != nil {
		return err
	}

	m.data = append(m.data, &routes{
		path:        path,
		method:      method,
		fn:          fn,
		middlewares: middlewares,
	})

	return nil
}

func (m *router) get(method, path string) *routes {
	for _, v := range m.data {
		if v.method == method && v.path == path {
			return v
		}
	}
	return nil
}

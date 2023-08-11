package httpd

import (
	"github.com/harley9293/nebulus/pkg/def"
	"net/http"
	"reflect"
)

type Context struct {
	Session def.Session

	r       *http.Request
	w       http.ResponseWriter
	service *Service

	status int
	err    error

	in  reflect.Value
	out any

	index       int
	middlewares []MiddlewareFunc
	handler     reflect.Value
}

func (c *Context) CreateSession(key string) {
	c.Session = c.service.NewSession(key)
}

func (c *Context) Next() {
	if c.index >= len(c.middlewares) {
		var params []reflect.Value
		if c.in.IsValid() {
			params = append(params, c.in)
		}
		params = append(params, reflect.ValueOf(c))
		result := c.handler.Call(params)

		if len(result) == 0 {
			c.w.Header().Set("Content-Type", "text/plain")
			c.out = []byte("")
		} else if result[0].Kind() == reflect.String {
			c.w.Header().Set("Content-Type", "text/plain")
			c.out = []byte(result[0].String())
		} else {
			c.w.Header().Set("Content-Type", "application/json")
			c.out = result[0].Interface()
		}
		return
	}

	c.index++
	c.middlewares[c.index-1](c)
}

func (c *Context) Error(status int, err error) {
	c.status = status
	c.err = err
}

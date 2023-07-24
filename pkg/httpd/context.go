package httpd

import (
	"net/http"
	"reflect"
)

type Context struct {
	Session *Session

	r       *http.Request
	w       http.ResponseWriter
	service *Service
	status  int

	in  reflect.Value
	out any

	index       int
	middlewares []MiddlewareFunc
	handler     reflect.Value
}

func (c *Context) CreateSession(key string) {
	c.Session = c.service.sm.new(key)
}

func (c *Context) Next() {
	if c.status != http.StatusOK {
		return
	}

	if c.index >= len(c.middlewares) {
		result := c.handler.Call([]reflect.Value{c.in, reflect.ValueOf(c)})
		c.out = result[0].Interface()
		return
	}

	c.index++
	c.middlewares[c.index-1](c)
}

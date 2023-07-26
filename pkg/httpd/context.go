package httpd

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type Context struct {
	Session *Session

	r       *http.Request
	w       http.ResponseWriter
	service *Service

	status int
	err    error

	in  reflect.Value
	out []byte

	index       int
	middlewares []MiddlewareFunc
	handler     reflect.Value
}

func (c *Context) CreateSession(key string) {
	c.Session = c.service.sm.new(key)
}

func (c *Context) Next() {
	if c.index >= len(c.middlewares) {
		result := c.handler.Call([]reflect.Value{c.in, reflect.ValueOf(c)})
		var err error
		c.out, err = json.Marshal(result[0].Interface())
		if err != nil {
			c.Error(http.StatusInternalServerError, err)
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

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
	status  int

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
	if c.status != http.StatusOK {
		return
	}

	if c.index >= len(c.middlewares) {
		result := c.handler.Call([]reflect.Value{c.in, reflect.ValueOf(c)})
		var err error
		c.out, err = json.Marshal(result[0].Interface())
		if err != nil {
			c.status = http.StatusInternalServerError
		}
		return
	}

	c.index++
	c.middlewares[c.index-1](c)
}

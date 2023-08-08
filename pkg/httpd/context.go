package httpd

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type Context struct {
	Session Session

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
		var params []reflect.Value
		if c.in.IsValid() {
			params = append(params, c.in)
		}
		params = append(params, reflect.ValueOf(c))
		result := c.handler.Call(params)

		if len(result) == 0 {
			c.w.Header().Set("Content-Type", "text/plain")
		} else if result[0].Kind() == reflect.String {
			c.w.Header().Set("Content-Type", "text/plain")
			c.out = []byte(result[0].String())
		} else {
			c.w.Header().Set("Content-Type", "application/json")
			var err error
			c.out, err = json.Marshal(result[0].Interface())
			if err != nil {
				c.Error(http.StatusInternalServerError, err)
			}
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

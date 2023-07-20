package httpd

import (
	"net/http"
	"reflect"
)

type Context struct {
	Session *Session

	r      *http.Request
	w      http.ResponseWriter
	sm     *sessionMng
	index  int
	h      *handlerData
	values []reflect.Value
	status int
	rsp    any
}

func (c *Context) CreateSession(key string) {
	c.Session = c.sm.new(key)
}

func (c *Context) Next() {
	if c.status != http.StatusOK {
		return
	}

	if c.index >= len(c.h.middlewares) {
		result := c.h.handler.Call(c.values)
		c.rsp = result[0].Interface()
		return
	}

	c.h.middlewares[c.index](c)
	c.index++
}

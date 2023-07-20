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

func (c *Context) Next() error {
	if c.status != http.StatusOK {
		return nil
	}

	if c.index >= len(c.h.middlewares) {
		result := c.h.handler.Call(c.values)
		c.rsp = result[0].Interface()
		return nil
	}

	err := c.h.middlewares[c.index](c)
	if err != nil {
		return err
	}
	c.index++
	return nil
}

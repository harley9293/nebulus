package httpd

import (
	"encoding/json"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"reflect"
)

type Context struct {
	Session *Session

	r      *http.Request
	w      http.ResponseWriter
	super  *handlerMng
	index  int
	h      *handlerData
	values []reflect.Value
}

func (c *Context) CreateSession(key string) {
	c.Session = c.super.sm.new(key)
}

func (c *Context) Next() error {
	if c.index >= len(c.h.middlewares) {
		result := c.h.handler.Call(c.values)
		err := json.NewEncoder(c.w).Encode(result[0].Interface())
		if err != nil {
			http.Error(c.w, "Internal Server Error", http.StatusInternalServerError)
			return errors.New("Internal Server Error")
		}

		return nil
	}

	err := c.h.middlewares[c.index](c)
	if err != nil {
		return err
	}
	c.index++
	return nil
}

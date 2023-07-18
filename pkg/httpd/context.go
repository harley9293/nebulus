package httpd

import (
	"encoding/json"
	log "github.com/harley9293/blotlog"
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
			log.Error("url: %s, err: %s", c.r.URL, err.Error())
			return err
		}

		log.Debug("url: %s, rsp: %+v", c.r.URL, result[0].Interface())
		return nil
	}

	err := c.h.middlewares[c.index](c)
	if err != nil {
		return err
	}
	c.index++
	return nil
}

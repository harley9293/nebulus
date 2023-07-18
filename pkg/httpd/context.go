package httpd

import "net/http"

type Context struct {
	Session *Session

	r     *http.Request
	super *handlerMng
}

func (c *Context) CreateSession(key string) {
	c.Session = c.super.sm.new(key)
}

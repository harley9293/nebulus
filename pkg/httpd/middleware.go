package httpd

import (
	"encoding/json"
	log "github.com/harley9293/blotlog"
	"net/http"
	"reflect"
)

func AuthMW(ctx *Context) {
	sessionCookie, err := ctx.r.Cookie("token")
	if err == nil {
		ctx.Session = ctx.service.sm.get(sessionCookie.Value)
	}

	if ctx.Session == nil {
		ctx.status = http.StatusUnauthorized
		return
	}

	ctx.Next()
}

func CookieMW(ctx *Context) {
	ctx.Next()

	if ctx.Session != nil {
		http.SetCookie(ctx.w, &http.Cookie{
			Name:  "token",
			Value: ctx.Session.id,
		})
	}
}

func defaultMW(ctx *Context) {
	h, ok := ctx.service.hm.data[ctx.r.URL.Path]
	if !ok {
		ctx.status = http.StatusNotFound
		return
	}

	if ctx.r.Method != h.method {
		ctx.status = http.StatusMethodNotAllowed
		return
	}

	arg := reflect.New(h.handler.Type().In(0))
	err := json.NewDecoder(ctx.r.Body).Decode(arg.Interface())
	if err != nil {
		ctx.status = http.StatusBadRequest
		return
	}
	ctx.in = arg.Elem()
	ctx.handler = h.handler
	if len(h.middlewares) > 0 {
		ctx.middlewares = append(ctx.middlewares, h.middlewares...)
	}

	log.Debug("url: %s, req: %+v", ctx.r.URL, arg.Elem().Interface())
	ctx.Next()

	if ctx.status != http.StatusOK {
		http.Error(ctx.w, http.StatusText(ctx.status), ctx.status)
		log.Error("url: %s err, statue: %d", ctx.r.URL, ctx.status)
		return
	}

	err = json.NewEncoder(ctx.w).Encode(ctx.out)
	if err != nil {
		ctx.status = http.StatusInternalServerError
		return
	}
	log.Debug("url: %s, rsp: %+v", ctx.r.URL, ctx.out)
}

package httpd

import (
	"encoding/json"
	log "github.com/harley9293/blotlog"
	"github.com/harley9293/nebulus/pkg/errors"
	"net/http"
	"reflect"
)

func AuthMW(ctx *Context) {
	sessionCookie, err := ctx.r.Cookie("token")
	if err == nil {
		ctx.Session = ctx.service.sm.get(sessionCookie.Value)
	}

	if ctx.Session == nil {
		ctx.Error(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
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

func CorsMW(ctx *Context) {
	ctx.w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	ctx.Next()
}

func RspPackMW(ctx *Context) {
	ctx.Next()

	rsp := make(map[string]interface{})
	if ctx.status == http.StatusOK {
		rsp["code"] = 0
		rsp["msg"] = "success"
		rsp["data"] = ctx.out
	} else {
		rsp["code"] = ctx.status
		rsp["msg"] = http.StatusText(ctx.status)
	}
	var err error
	ctx.out, err = json.Marshal(rsp)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err)
	}
}

func LogMW(ctx *Context) {
	log.Info("url: %s, req: %+v", ctx.r.URL, ctx.in.Interface())
	ctx.Next()

	if ctx.status != http.StatusOK {
		log.Error("url: %s, err, statue: %d, msg: %s", ctx.r.URL, ctx.status, ctx.err.Error())
		return
	}
	log.Info("url: %s, rsp: %s", ctx.r.URL, string(ctx.out))
}

func baseMW(ctx *Context) {
	h, ok := ctx.service.hm.data[ctx.r.URL.Path]
	if !ok {
		ctx.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}

	if ctx.r.Method != h.method {
		ctx.Error(http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	arg := reflect.New(h.handler.Type().In(0))
	err := json.NewDecoder(ctx.r.Body).Decode(arg.Interface())
	if err != nil {
		ctx.Error(http.StatusBadRequest, err)
		return
	}
	ctx.in = arg.Elem()
	ctx.handler = h.handler
	if len(h.middlewares) > 0 {
		ctx.middlewares = append(ctx.middlewares, h.middlewares...)
	}

	ctx.Next()

	if ctx.status == http.StatusOK {
		_, err = ctx.w.Write(ctx.out)
		if err != nil {
			log.Error("ctx.w.Write() failed, err:%s", err.Error())
			return
		}
	} else {
		http.Error(ctx.w, http.StatusText(ctx.status), ctx.status)
		return
	}
}

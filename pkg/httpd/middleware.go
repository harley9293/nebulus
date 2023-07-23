package httpd

import "net/http"

func AuthMW(ctx *Context) {
	if ctx.Session == nil {
		ctx.status = http.StatusUnauthorized
		return
	}

	ctx.Next()
}

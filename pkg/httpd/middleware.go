package httpd

import "net/http"

func SessionMiddleware(ctx *Context) {
	sessionCookie, err := ctx.r.Cookie("session_id")
	if err != nil {
		ctx.status = http.StatusUnauthorized
		return
	}

	if ctx.sm.get(sessionCookie.Value) == nil {
		ctx.status = http.StatusUnauthorized
		return
	}

	ctx.Next()
}

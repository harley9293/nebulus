package httpd

import "net/http"

func SessionMiddleware(ctx *Context) error {
	_, err := ctx.r.Cookie("session_id")
	if err != nil {
		http.Error(ctx.w, "Unauthorized", http.StatusUnauthorized)
		return err
	}

	return ctx.Next()
}

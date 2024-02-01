package middleware

import (
	"context"
	"net/http"
)

type ContextBundle struct {
	Key int
	Val any
}

func Context(next http.Handler, ctxBundles ...*ContextBundle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, ctxBundle := range ctxBundles {
			ctx = context.WithValue(ctx, ctxBundle.Key, ctxBundle.Val)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

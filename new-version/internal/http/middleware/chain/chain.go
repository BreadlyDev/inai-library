package chain

import (
	"context"
	"net/http"
)

type Middleware func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool

func Chain(
	ctx context.Context,
	handler func(w http.ResponseWriter, r *http.Request),
	middlewares ...Middleware,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, mw := range middlewares {
			ok := mw(ctx, w, r)
			if !ok {
				return
			}
		}

		handler(w, r)
	})
}

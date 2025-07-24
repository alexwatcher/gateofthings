package middlewares

import "net/http"

func Group(mux *http.ServeMux, prefix string, register func(mux *http.ServeMux), middlewares ...Middleware) {
	subMux := http.NewServeMux()
	h := http.Handler(subMux)
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	mux.Handle(prefix+"/", http.StripPrefix(prefix, h))
	register(subMux)
}

func GroupMiddlewares(middlewares ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for _, middleware := range middlewares {
			h = middleware(h)
		}
		return h
	}
}

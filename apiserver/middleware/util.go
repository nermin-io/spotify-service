package middleware

import "net/http"

// Apply wraps a list of middleware around the http.Handler.
// Please note that requests flow through the middleware in the order specified,
// for example the second middleware can see effects of the first middleware, and so on.
func Apply(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

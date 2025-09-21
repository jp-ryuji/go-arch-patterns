package http

import (
	"net/http"
)

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// ChainMiddleware chains multiple middleware functions
func ChainMiddleware(mw ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			head := final
			for i := len(mw) - 1; i >= 0; i-- {
				head = mw[i](head)
			}
			head.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a production system, you would use a proper logger
		// For now, we'll just print to stdout
		println("Request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Digivate-Labs-Pvt-Ltd/dvlutil"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		next = methodNotAllowed(next)
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (wr *wrappedWriter) Write(b []byte) (int, error) {

	if wr.statusCode == http.StatusMethodNotAllowed {
		b, _ := json.Marshal(dvlutil.Response{
			Status: dvlutil.StatusCodeNotOK,
			Msg:    "Method not allowed",
		})

		wr.Header().Set("Content-Type", "application/json")
		return wr.ResponseWriter.Write(b)
	}

	return wr.ResponseWriter.Write(b)
}

func methodNotAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(wrapped, r)
	})
}

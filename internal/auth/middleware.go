package auth

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// Wrap response writer
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Execute actual handler
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		fmt.Printf(
			"%s %s | Status: %d | Size: %d bytes | Time: %v\n",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			rw.size,
			duration,
		)
	})
}

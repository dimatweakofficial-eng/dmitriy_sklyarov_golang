package middleware

import "net/http"

type WraperWrite struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WraperWrite) WriteHeader(statusCode int) {
	w.WriteHeader(statusCode)
}

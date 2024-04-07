package server

import (
	"io"
	"net/http"
)

// Logger handler
type ResponceLogWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rlw *ResponceLogWriter) Write(b []byte) (int, error) {
	size, err := rlw.ResponseWriter.Write(b)
	rlw.size = size
	return size, err
}

func (rlw *ResponceLogWriter) WriteHeader(statusCode int) {
	rlw.ResponseWriter.WriteHeader(statusCode)
	rlw.statusCode = statusCode
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

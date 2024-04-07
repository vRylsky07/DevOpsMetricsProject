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

// gzip

type gzipWriter struct {
	http.ResponseWriter
	Writer       io.Writer
	AllowedTypes *[]string
}

func (w gzipWriter) Write(b []byte) (int, error) {

	header := w.Header().Get("Content-Type")

	if w.AllowedTypes != nil {
		for _, s := range *w.AllowedTypes {
			if s == header {
				w.Header().Set("Content-Encoding", "gzip")
				return w.Writer.Write(b)
			}
		}
	}

	return w.ResponseWriter.Write(b)
}

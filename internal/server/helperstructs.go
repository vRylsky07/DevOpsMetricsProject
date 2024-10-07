package server

import (
	"compress/gzip"
	"net/http"
	"strings"
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
	Writer       *gzip.Writer
	AllowedTypes *[]string
}

func (w gzipWriter) WriteHeader(statusCode int) {

	header := w.Header().Get("Content-Type")

	if w.AllowedTypes != nil && statusCode < 300 {
		for _, s := range *w.AllowedTypes {
			if strings.Contains(header, s) {
				w.Header().Set("Content-Encoding", "gzip")
			}
		}
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w gzipWriter) Write(b []byte) (int, error) {

	if w.Header().Get("Content-Encoding") == "gzip" {
		defer w.Writer.Close()
		return w.Writer.Write(b)
	}

	return w.ResponseWriter.Write(b)
}

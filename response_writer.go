package fire

import (
	"net/http"
)

type FireResponseWriter struct {
	Code           int
	responseWriter http.ResponseWriter
	wroteHeader    bool
}

func NewFireResponseWriter(rw http.ResponseWriter) *FireResponseWriter {
	return &FireResponseWriter{responseWriter: rw, wroteHeader: false}
}

func (f *FireResponseWriter) Header() http.Header {
	return f.responseWriter.Header()
}

func (f *FireResponseWriter) Write(buf []byte) (int, error) {
	if !f.wroteHeader {
		f.WriteHeader(http.StatusOK)
	}
	return f.responseWriter.Write(buf)
}

func (f *FireResponseWriter) WriteHeader(code int) {
	f.Code = code
	f.wroteHeader = true
	f.responseWriter.WriteHeader(code)
}

package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
	"xrf197ilz35aq0/internal"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer // // A buffer for the response body
}

// WriteHeader writes the status code to the response.
func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Header returns the headers of the underlying response writer.
func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Write to the buffer
	return w.ResponseWriter.Write(b)
}

// LoggerHandler is a middleware that logs requests.
type LoggerHandler struct {
	logger internal.Logger
}

func (lh *LoggerHandler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := internal.GenerateRequestId()
		start := time.Now()

		logPrefix := fmt.Sprintf("requestId='%s'", requestId)
		lh.logger.SetPrefix(logPrefix)

		lh.logger.Info(fmt.Sprintf("event=request :: method=%s :: url=%s :: remoteAddr=%s :: userAgent=%s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent()))

		// Wrap the response writer to capture the status code.
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		wrappedWriter.Header().Set("Request-Trace-Id", requestId)

		// Call the next handler.
		next.ServeHTTP(wrappedWriter, r)

		// Stop the timer.
		duration := time.Since(start)

		if wrappedWriter.status >= 400 {
			lh.logger.Error(fmt.Sprintf("event=response :: success=false :: url=%s :: status=%d :: duration=%dms error=%s",
				r.URL.Path,
				wrappedWriter.status,
				int(duration.Milliseconds()),
				wrappedWriter.body.String()))
		} else {
			lh.logger.Info(fmt.Sprintf("event=response :: success=true :: url=%s :: status=%d :: duration=%dms",
				r.URL.Path,
				wrappedWriter.status,
				int(duration.Milliseconds())))
		}
	})
}

func NewLoggerHandler(logger internal.Logger) *LoggerHandler {
	return &LoggerHandler{
		logger: logger,
	}
}

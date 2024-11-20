package middleware

import (
	"fmt"
	"net/http"
	"time"
	"xrf197ilz35aq0"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader writes the status code to the response.
func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// LoggerHandler is a middleware that logs requests.
type LoggerHandler struct {
	logger xrf197ilz35aq0.Logger
}

func (lh *LoggerHandler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

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

		// Call the next handler.
		next.ServeHTTP(wrappedWriter, r)

		// Stop the timer.
		duration := time.Since(start)

		lh.logger.Info(fmt.Sprintf("event=response :: status=%d duration=%s", wrappedWriter.status, duration))
	})
}

func NewLoggerHandler(logger xrf197ilz35aq0.Logger) *LoggerHandler {
	return &LoggerHandler{
		logger: logger,
	}
}
package middleware

import (
	"fmt"
	"github.com/nermin-io/spotify-service/logging"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		start := time.Now()
		lrw := loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(&lrw, r)
		duration := time.Since(start)
		logMessage := fmt.Sprintf("%s %s (%s)", r.Method, r.URL.Path, duration)
		logger.Info(logMessage,
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Int("status", lrw.status),
			zap.Int("size", lrw.size),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}

// loggingResponseWriter is a wrapper for http.ResponseWriter to capture status code and size.
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

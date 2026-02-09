package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func LoggingMiddleWare(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &ResponseWriterWrapper{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			logger.Info("request started", zap.String("method", r.Method), zap.String("url", r.URL.Path))

			defer func() {
				duration := time.Since(start)
				logger.Info("request finished", zap.Int("status", rw.status), zap.Duration("duration", duration))
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

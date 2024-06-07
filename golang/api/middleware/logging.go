package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

const (
	logLevel = "INFO"
)

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingWriter(w http.ResponseWriter) *loggingWriter {
	return &loggingWriter{
		ResponseWriter: w,
		statusCode:     http.StatusInternalServerError,
	}
}

func (lw *loggingWriter) WriteHeader(statusCode int) {
	lw.statusCode = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		spanCtx := trace.SpanContextFromContext(r.Context())

		rlw := newLoggingWriter(w)

		initTime := time.Now()
		next.ServeHTTP(rlw, r)
		duration := time.Since(initTime)

		defer func() {
			if rlw.statusCode >= 500 && rlw.statusCode < 600 {
				slog.Error("", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.Duration("Duration", duration), slog.String("method", r.Method), slog.String("uri", r.RequestURI), slog.String("trace-id", spanCtx.TraceID().String()))
			} else if rlw.statusCode >= 400 && rlw.statusCode < 500 {
				slog.Warn("", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.Duration("Duration", duration), slog.String("method", r.Method), slog.String("uri", r.RequestURI), slog.String("trace-id", spanCtx.TraceID().String()))
			} else {
				slog.Info("", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.Duration("Duration", duration), slog.String("method", r.Method), slog.String("uri", r.RequestURI), slog.String("trace-id", spanCtx.TraceID().String()))
			}
		}()
	})
}

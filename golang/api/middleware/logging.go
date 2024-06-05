package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

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
		// Tracerが仕込まれていた場合, TracerIDを記述する
		if spanCtx.HasTraceID() {
			slog.Info("request", slog.String("method", r.Method), slog.String("uri", r.RequestURI), slog.String("trace-id", spanCtx.TraceID().String()))
		}

		rlw := newLoggingWriter(w)

		next.ServeHTTP(rlw, r)

		defer func() {
			if rlw.statusCode >= 500 && rlw.statusCode < 600 {
				slog.Error("Response", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.String("trace-id", spanCtx.TraceID().String()))
			} else if rlw.statusCode >= 400 && rlw.statusCode < 500 {
				slog.Warn("Response", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.String("trace-id", spanCtx.TraceID().String()))
			} else {
				slog.Info("Response", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.String("trace-id", spanCtx.TraceID().String()))
			}
		}()
	})
}

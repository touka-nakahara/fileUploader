package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
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

func LoggingMiddleware(next http.Handler, httpLogger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ここで個別ログを作っていいものか...あとあとファイルオープンもここですることになるし

		//RV RemoteAddrを入れるべきかはわからないが...一意にする必要があった
		httpLogger.Info("Rquest", slog.String("method", r.Method), slog.String("uri", r.RequestURI), slog.String("RemoteAddr", r.RemoteAddr))

		rlw := newLoggingWriter(w)

		next.ServeHTTP(rlw, r)

		defer func() {
			httpLogger.Info("Response", slog.String("status", fmt.Sprintf("%d", rlw.statusCode)), slog.String("RemoteAddr", r.RemoteAddr))
		}()
	})
}

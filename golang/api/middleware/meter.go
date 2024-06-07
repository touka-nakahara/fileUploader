package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type httpMetricMiddleware struct {
	next                     http.Handler
	requestDurationHistogram metric.Int64Histogram
	counter                  metric.Int64Counter
}

func NewMetricMiddleware(meter metric.Meter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		durationHistogram, _ := meter.Int64Histogram("http.server.latency", metric.WithUnit("ms"))
		counter, _ := meter.Int64Counter("access.counter", metric.WithDescription("Number of API calls"), metric.WithUnit("{call}"))
		return &httpMetricMiddleware{
			next:                     next,
			requestDurationHistogram: durationHistogram,
			counter:                  counter,
		}
	}
}

func (h *httpMetricMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rlw := newStatusCodeCapturerWriter(w)

	// ここにInt64Histogramで記録する
	initTime := time.Now()
	h.next.ServeHTTP(rlw, r)
	duration := time.Since(initTime)

	// traceID := trace.SpanContextFromContext(r.Context()).TraceID().String()

	// 色々セット
	//TODO TraceID乗っける
	//RV nakahara TraceIDを乗っけるとユニークになるので統計情報としては機能しなくなる
	metricAttributes := attribute.NewSet(
		semconv.HTTPURL(r.RequestURI),
		semconv.HTTPRequestContentLength(int(r.ContentLength)),
		semconv.HTTPMethod(r.Method),
		semconv.HTTPStatusCode(rlw.statusCode),
		// attribute.String("traceID", traceID),
	)

	// メトリクス送信
	h.requestDurationHistogram.Record(
		r.Context(),
		duration.Microseconds(),
		metric.WithAttributeSet(metricAttributes),
	)

	// カウンタ
	h.counter.Add(
		r.Context(),
		1,
		// metric.WithAttributeSet(attribute.NewSet(attribute.String("traceID", traceID))),
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newStatusCodeCapturerWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

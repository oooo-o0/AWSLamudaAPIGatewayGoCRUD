package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// レスポンスのステータスを取得するためのラッパー
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf("Completed %d %s in %v", lrw.statusCode, http.StatusText(lrw.statusCode), duration)
	})
}

// ステータスコード取得用ラッパー
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

package logger

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	bytes int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

func LoggerMiddleware(log *slog.Logger) func(http.Handler) http.Handler {

	middleware := func(next http.Handler) http.Handler {

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := &responseWriter{ResponseWriter: w}

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			tnow := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.String("duration", time.Since(tnow).String()),
					slog.Int("bytes", ww.bytes),
				)
			}()

			next.ServeHTTP(ww, r)
		})

		return handler
	}

	return middleware
}

func Logger(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	log, ok := ctx.Value("logger").(*slog.Logger)
	if !ok || log == nil {
		http.Error(w, "internal server error (logger missing)", http.StatusInternalServerError)
		return false
	}
	ww := &responseWriter{ResponseWriter: w}

	entry := log.With(
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
	)

	tnow := time.Now()
	entry.Info("request started")

	defer func() {
		entry.Info("request completed",
			slog.String("duration", time.Since(tnow).String()),
			slog.Int("bytes", ww.bytes),
		)
	}()

	return true
}

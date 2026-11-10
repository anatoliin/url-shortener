package router

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				if rec := recover(); rec != nil && rec != http.ErrAbortHandler {
					logger.Error().
						Interface("recovered", rec).
						Bytes("stack", debug.Stack()).
						Msg("incoming request panic")
					ww.WriteHeader(http.StatusInternalServerError)
				}
				logger.Info().
					Str("uri", r.RequestURI).
					Str("method", r.Method).
					Int("status_code", ww.Status()).
					Int("bytes_out", ww.BytesWritten()).
					Dur("duration", time.Since(start)).
					Msg("request processed")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

package router

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				if r := recover(); r != nil && r != http.ErrAbortHandler {
					logger.Error().Interface("recovered", r).Bytes("stack", debug.Stack()).Msg("incoming request panic")
					ww.WriteHeader(http.StatusInternalServerError)
				}
				logger.Info().Fields(map[string]any{
					"uri":         r.RequestURI,
					"method":      r.Method,
					"status_code": ww.Status(),
					"bytes_out":   ww.BytesWritten(),
					"duration":    time.Since(start),
				}).Msg("request processed")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

package logginghandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetUUID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

func Logger(r *http.Request) zerolog.Logger {
	logger := log.With().Str("uuid", GetUUID(r)).Logger()

	return logger
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New().String()
		r.Header.Set("X-Request-ID", uuid)
		logger := Logger(r)
		logger.Info().
			Str("uuid", uuid).
			Str("method", r.Method).
			Str("user-agent", r.UserAgent()).
			Str("proto", r.Proto).
			Str("referer", r.Referer()).
			Str("request-url", r.URL.String()).
			Str("remote", r.RemoteAddr).
			Msg("")

		w.Header().Set("X-Request-ID", uuid)
		next.ServeHTTP(w, r)
	})
}

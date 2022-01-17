// Package logginghandler is a simple, zerolog based, request logging http middleware.
// It also sets `X-Request-ID` in the request and response headers.
package logginghandler

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// GetUUID gets the requests UUID from a request.
func GetUUID(r *http.Request) string {
	uuid, ok := hlog.IDFromRequest(r)
	if !ok {
		return ""
	}

	return uuid.String()
}

// Logger returns a logger with the UUID set.
func Logger(r *http.Request) zerolog.Logger {
	return *hlog.FromRequest(r)
}

func Handler(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		chain := alice.New(
			hlog.NewHandler(log),
			hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Str("proto", r.Proto).
					Stringer("request-url", r.URL).
					Int("status", status).
					Int("size", size).
					Dur("duration", duration).
					Msg("")
			}),
			hlog.RemoteAddrHandler("remote"),
			hlog.UserAgentHandler("user-agent"),
			hlog.RefererHandler("referer"),
			hlog.RequestIDHandler("uuid", "X-Request-ID"),
		).Then(next)

		return chain
	}
}

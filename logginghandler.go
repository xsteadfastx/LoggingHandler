// Package logginghandler is a simple, zerolog based, request logging http middleware.
// It also sets `X-Request-ID` in the request and response headers.
package logginghandler

import (
	"context"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// GetUUID gets the requests UUID from a request.
func GetUUID(r *http.Request) (string, bool) {
	uuid, ok := hlog.IDFromRequest(r)
	if !ok {
		return "", false
	}

	return uuid.String(), true
}

// FromRequest returns a logger with the UUID set from request.
// If no one could be found, it will return the global one.
func FromRequest(r *http.Request) zerolog.Logger {
	l := hlog.FromRequest(r)

	if l.GetLevel() == zerolog.Disabled {
		return log.Logger
	}

	return *hlog.FromRequest(r)
}

// FromCtx returns a logger with the UUID set from ctx.
// If no one could be found, it will return the global one.
func FromCtx(ctx context.Context) zerolog.Logger {
	l := *log.Ctx(ctx)

	if l.GetLevel() == zerolog.Disabled {
		return log.Logger
	}

	return l
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

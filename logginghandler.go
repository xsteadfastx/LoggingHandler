// Package logginghandler is a simple, zerolog based, request logging http middleware.
// It also sets `X-Request-ID` in the request and response headers.
package logginghandler

import (
	"context"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

const (
	UUIDKey    = "uuid"
	UUIDHeader = "X-Request-ID"
)

func init() { //nolint:gochecknoinits
	zerolog.DefaultContextLogger = &log.Logger
}

// GetUUID gets the requests UUID from a request.
func GetUUID(r *http.Request) (string, bool) {
	uuid, ok := hlog.IDFromRequest(r)
	if !ok {
		return "", false
	}

	return uuid.String(), true
}

// FromRequest returns a logger set from request.
// If no one could be found, it will return the global one.
func FromRequest(r *http.Request) *zerolog.Logger {
	return hlog.FromRequest(r)
}

// FromCtx returns a logger set from ctx.
// If no one could be found, it will return the global one.
func FromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// RequestIDHandler looks in the header for an existing request id. Else it will create one.
func RequestIDHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(UUIDHeader)

			if id != "" {
				ctx := r.Context()

				log := zerolog.Ctx(ctx)

				uuid, err := xid.FromString(id)
				if err != nil {
					log.Error().Err(err).Msg("couldnt parse uuid")

					hlog.RequestIDHandler(UUIDKey, UUIDHeader)(next).ServeHTTP(w, r)

					return
				}

				ctx = hlog.CtxWithID(ctx, uuid)
				r = r.WithContext(ctx)

				log.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str(UUIDKey, uuid.String())
				})

				w.Header().Set(UUIDHeader, uuid.String())

				next.ServeHTTP(w, r)
			} else {
				hlog.RequestIDHandler(UUIDKey, UUIDHeader)(next).ServeHTTP(w, r)
			}
		})
	}
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
			RequestIDHandler(),
		).Then(next)

		return chain
	}
}

// Package logginghandler is a simple, zerolog based, request logging http middleware.
// It also sets `X-Request-ID` in the request and response headers.
package logginghandler

import (
	"net/http"

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

// Handler sets up all the logging.
func Handler(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return hlog.NewHandler(logger)(
			hlog.RemoteAddrHandler("remote")(
				hlog.UserAgentHandler("user-agent")(
					hlog.RefererHandler("referer")(
						hlog.MethodHandler("method")(
							hlog.RequestIDHandler("uuid", "X-Request-ID")(
								hlog.URLHandler("request-url")(
									next,
								),
							),
						),
					),
				),
			),
		)
	}
}

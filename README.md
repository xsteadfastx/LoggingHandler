# logginghandler

[![Build Status](https://ci.xsfx.dev/api/badges/xsteadfastx/logginghandler/status.svg)](https://ci.xsfx.dev/xsteadfastx/logginghandler)

Just a simple zerolog based request logging http middleware. It also sets a `X-Request-ID` in the request and response headers.

## Install

        go get -v git.xsfx.dev/xsteadfastx/logginghandler

## Usage

        handler := logginghandler.Handler(http.HandlerFunc(myHandler))
        http.Handle("/", handler)
        log.Fatal().Msg(http.ListenAndServe(":5000", nil).Error())

In other handlers you can access the UUID:

        func anotherHandler(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "your uuid is: %s", logginghandler.GetUUID(r))
        }

The already prepared logger is also available:

        l := loggerhandler.Logger(r)
        l.Info().Msg("foo bar")

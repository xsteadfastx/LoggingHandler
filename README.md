# logginghandler

[![Build Status](https://ci.xsfx.dev/api/badges/xsteadfastx/logginghandler/status.svg)](https://ci.xsfx.dev/xsteadfastx/logginghandler)
[![Go Reference](https://pkg.go.dev/badge/go.xsfx.dev/logginghandler.svg)](https://pkg.go.dev/go.xsfx.dev/logginghandler)
[![Go Report Card](https://goreportcard.com/badge/go.xsfx.dev/logginghandler)](https://goreportcard.com/report/go.xsfx.dev/logginghandler)

Just a simple zerolog based request logging http middleware. It also sets a `X-Request-ID` in the request and response headers.

Powered by [github.com/rs/zerolog/hlog](https://github.com/rs/zerolog) and [github.com/justinas/alice](https://github.com/justinas/alice).

## Install

```shell
go get -v go.xsfx.dev/logginghandler
```

## Usage

```golang
logger := log.With().Logger()

handler := logginghandler.Handler(logger)(
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        log := logginghander.Logger(r)
        log.Info().Msg("hello world")

        w.WriteHeader(http.StatusOK)

        return
    })
)

http.Handle("/", handler)
log.Fatal().Msg(http.ListenAndServe(":5000", nil).Error())
```

or with [alice](https://github.com/justinas/alice)

```golang
logger := log.With().Logger()
chain := alice.New(logginghandler.Handler(logger)).Then(
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        log := logginghander.Logger(r)
        log.Info().Msg("hello world")

        w.WriteHeader(http.StatusOK)

        return
    })
)

http.Handle("/", chain)

log.Fatal().Err(http.ListenAndServe(":5000", nil)).Msg("goodbye")
```

In other handlers you can access the UUID:

```golang
func anotherHandler(w http.ResponseWriter, r *http.Request) {
    log := logginghandler.FromRequest(r)

    uuid, ok := logginghandler.GetUUID(r)
    if !ok {
        log.Error().Err(err).Msg("could not find uuid")
        w.WriteHeader(http.StatusInternalServerError)

        return
    }

    fmt.Fprintf(w, "your uuid is: %s", uuid)

    return
}
```

The already prepared logger is also available:

```golang
l := logginghandler.FromRequest(r)
l.Info().Msg("foo bar")
```

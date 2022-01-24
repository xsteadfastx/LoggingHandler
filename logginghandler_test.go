package logginghandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"go.xsfx.dev/logginghandler"
)

func Example() {
	logger := log.With().Logger()

	handler := logginghandler.Handler(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logginghandler.FromRequest(r)

		logger.Info().Msg("this is a request")

		w.WriteHeader(http.StatusOK)
	}))

	http.Handle("/", handler)
	log.Fatal().Msg(http.ListenAndServe(":5000", nil).Error())
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestUUID(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/test", nil)
	assert.NoError(err)

	rr := httptest.NewRecorder()
	handler := logginghandler.Handler(log.With().Logger())(http.HandlerFunc(testHandler))

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(rr.Header().Get("X-Request-ID"))
}

func TestFromCtx(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/test", nil)
	assert.NoError(err)

	// Create buffer to store output.
	var output bytes.Buffer

	rr := httptest.NewRecorder()
	handler := logginghandler.Handler(zerolog.New(&output))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logginghandler.FromCtx(r.Context())
		log.Info().Msg("hello world")
	}))

	handler.ServeHTTP(rr, req)

	logs := strings.Split(output.String(), "\n")
	assert.Len(logs, 3)

	var jOut struct{ UUID string }

	err = json.Unmarshal([]byte(logs[0]), &jOut)
	assert.NoError(err)

	assert.NotEmpty(jOut)
}

package logginghandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"go.xsfx.dev/logginghandler"
)

func Example() {
	logger := log.With().Logger()

	handler := logginghandler.Handler(
		logger,
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logginghandler.FromRequest(r)

			logger.Info().Msg("this is a request")

			w.WriteHeader(http.StatusOK)
		}),
	)

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

	assert.NotEmpty(rr.Header().Get(logginghandler.UUIDHeader))
}

func TestFromCtx(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/test", nil)
	assert.NoError(err)

	// Create buffer to store output.
	var output bytes.Buffer

	rr := httptest.NewRecorder()
	handler := logginghandler.Handler(
		zerolog.New(&output),
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logginghandler.FromCtx(r.Context())
			log.Info().Msg("hello world")
		}),
	)

	handler.ServeHTTP(rr, req)

	logs := strings.Split(output.String(), "\n")
	assert.Len(logs, 3)

	var jOut struct{ UUID string }

	err = json.Unmarshal([]byte(logs[0]), &jOut)
	assert.NoError(err)

	assert.NotEmpty(jOut)
}

func TestRequestIDHandler(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	handler := logginghandler.RequestIDHandler()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := hlog.FromRequest(r)
			log.Info().Msg("hello from TestRequestID")
		}),
	)

	id := "cfrj1ro330reqgvfpgu0"

	// Create buffer to store output.
	var output bytes.Buffer

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	assert.NoError(err)

	rr := httptest.NewRecorder()

	h := hlog.NewHandler(zerolog.New(&output))(handler)

	h.ServeHTTP(rr, req)

	assert.NotEmpty(rr.Header().Get(logginghandler.UUIDHeader))
	assert.NotEqual(rr.Header().Get(logginghandler.UUIDHeader), id)

	// Now test with request id in header.

	nr := httptest.NewRecorder()

	nReq, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	assert.NoError(err)

	nReq.Header.Add(logginghandler.UUIDHeader, id)

	h.ServeHTTP(nr, nReq)

	assert.NotEmpty(nr.Header().Get(logginghandler.UUIDHeader))
	assert.Equal(nr.Header().Get(logginghandler.UUIDHeader), id)

	logs := strings.Split(output.String(), "\n")
	assert.Len(logs, 3)

	getUUID := func(l string) (string, error) {
		var out struct{ UUID string }

		err := json.Unmarshal([]byte(l), &out)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal log: %w", err)
		}

		return out.UUID, nil
	}

	uuid1, err := getUUID(logs[0])
	assert.NoError(err)

	assert.NotEqual(id, uuid1)

	uuid2, err := getUUID(logs[1])
	assert.NoError(err)
	assert.Equal(id, uuid2)
}

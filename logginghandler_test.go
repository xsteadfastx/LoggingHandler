package logginghandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.xsfx.dev/logginghandler"
)

func Example() {
	handler := logginghandler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	http.Handle("/", handler)
	log.Fatal().Msg(http.ListenAndServe(":5000", nil).Error())
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("got request")
}

func TestUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/test", nil)
	assert.NoError(err)

	rr := httptest.NewRecorder()
	handler := logginghandler.Handler(http.HandlerFunc(testHandler))

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(rr.Header().Get("X-Request-ID"))
	log.Print(rr.Header())
}

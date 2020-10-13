package logginghandler_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.xsfx.dev/xsteadfastx/logginghandler"
	"github.com/stretchr/testify/assert"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("got request")
}

func TestUUID(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/test", nil)
	assert.NoError(err)

	rr := httptest.NewRecorder()
	handler := logginghandler.Handler(http.HandlerFunc(testHandler))

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(rr.Header().Get("X-Request-ID"))
	log.Print(rr.Header())
}

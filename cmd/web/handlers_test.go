package main

import (
	"net/http"
	"testing"

	"github.com/Danvs60/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// application struct for end-to-end testing
	// add mock loggers (discard any log)
	app := newTestApplication(t)

	// test server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

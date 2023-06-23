package server_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/matryer/is"

	"canvas/integrationtest"
	"canvas/k"
)

func TestServer_Start(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("starts the server and listens for requests", func(t *testing.T) {
		is := is.New(t)

		cleanup := integrationtest.CreateServer()
		defer cleanup()

		healthEndpoint, err := url.JoinPath("http://localhost:8081/", k.GlobalPrefix, "health")
		is.NoErr(err)
		resp, err := http.Get(healthEndpoint)
		is.NoErr(err)
		is.Equal(http.StatusOK, resp.StatusCode)
	})
}

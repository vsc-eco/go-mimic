package e2e_test

import (
	"io"
	"mimic/modules/api"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHttpHandler(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		url := makeGoMimicEndpoint(goMimicPort)

		resp, err := http.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		buf, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, []byte(api.RootMsg), buf)
	})

	t.Run("healtcheck", func(t *testing.T) {
		url := makeGoMimicEndpoint(goMimicPort, "health")

		resp, err := http.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

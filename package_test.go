package requestid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	assert.Equal(t, defaultNamespace, Default())
	assert.Equal(t, defaultNamespace, Default(), "DefaultNamespace should return the same instance")

	ns1 := New(ResponseHeader("X-REQ-ID"))
	SetDefault(ns1)
	defer SetDefault(defaultNamespace)

	assert.Equal(t, ns1, Default())
	assert.Equal(t, ns1, Default(), "DefaultNamespace should return the same instance")
}

func TestDefaultSimpleWrap(t *testing.T) {
	backupDefaultNamespace := defaultNamespace
	defer func() { defaultNamespace = backupDefaultNamespace }()

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(
		Wrap(baseHandler),
	)
	defer ts.Close()

	t.Run("request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))
	})
}

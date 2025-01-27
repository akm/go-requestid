package requestid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	assert.Equal(t, degaultNamespace, DefaultNamespace())
	assert.Equal(t, degaultNamespace, DefaultNamespace(), "DefaultNamespace should return the same instance")

	ns1 := New(ResponseHeader("X-REQ-ID"))
	SetDefaultNamespace(ns1)
	defer SetDefaultNamespace(degaultNamespace)

	assert.Equal(t, ns1, DefaultNamespace())
	assert.Equal(t, ns1, DefaultNamespace(), "DefaultNamespace should return the same instance")
}

func TestDefaultSimpleWrap(t *testing.T) {
	backupDefaultNamespace := degaultNamespace
	defer func() { degaultNamespace = backupDefaultNamespace }()

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(
		Wrap(baseHandler),
	)
	defer ts.Close()

	t.Run("request", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))
	})
}

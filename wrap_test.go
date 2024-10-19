package requestid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(
		Wrap(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestIdOnHeader := r.Header.Get("X-Request-ID")
				requestIdFromContext := Get(r.Context())
				if requestIdOnHeader != "" {
					assert.Equal(t, requestIdOnHeader, requestIdFromContext)
				} else {
					assert.NotEmpty(t, requestIdFromContext)
					assert.Len(t, requestIdFromContext, 8)
				}
				baseHandler.ServeHTTP(w, r)
			}),
			RequestHeader("X-Request-ID"),
		),
	)
	defer ts.Close()

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL, nil)
		req.Header.Set("X-Request-ID", "in-header")
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "in-header", resp.Header.Get("X-Request-ID"))

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))
	})

	t.Run("request without X-Request-ID header", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))
	})
}

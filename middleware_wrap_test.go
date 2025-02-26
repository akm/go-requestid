package requestid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareWrap(t *testing.T) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})
	generatedCode := "genrated-code"
	mw := New(
		RequestHeader("X-Request-ID"),
		Generator(func() string { return generatedCode }),
	)
	ts := httptest.NewServer(
		mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestIDOnHeader := r.Header.Get("X-Request-ID")
			requestIDFromContext := mw.header.Get(r.Context())
			if requestIDOnHeader != "" {
				assert.Equal(t, requestIDOnHeader, requestIDFromContext)
			} else {
				assert.NotEmpty(t, requestIDFromContext)
				assert.Equal(t, generatedCode, requestIDFromContext)
			}
			baseHandler.ServeHTTP(w, r)
		})),
	)
	defer ts.Close()

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil) //nolint:noctx
		req.Header.Set("X-Request-ID", "in-header")
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "in-header", resp.Header.Get("X-Request-ID"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))
	})

	t.Run("request without X-Request-ID header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil) //nolint:noctx
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, generatedCode, resp.Header.Get("X-Request-ID"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))
	})
}

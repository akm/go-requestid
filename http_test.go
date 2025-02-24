package requestid

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestIdProviderWrapper(t *testing.T) {
	provider := func(*http.Request) string { return "generated" }
	wrapper := requestIDProviderWrapper(provider, "X-Request-ID")
	t.Run("request with X-Request-ID header", func(t *testing.T) {
		req := &http.Request{Header: http.Header{}}
		req.Header.Set("X-Request-ID", "in-header")
		assert.Equal(t, "in-header", wrapper(req))
	})
	t.Run("request without X-Request-ID header", func(t *testing.T) {
		assert.Equal(t, "generated", wrapper(new(http.Request)))
	})
}

func TestGeneratorProvider(t *testing.T) {
	generator := func() string {
		return "test"
	}
	provider := generatorProvider(generator)
	assert.Equal(t, "test", provider(nil))
}

func TestSelectNotEmpty(t *testing.T) {
	providerEmpty := func(*http.Request) string { return "" }
	provider := selectNotEmpty(providerEmpty)
	t.Run("empty", func(t *testing.T) {
		req := new(http.Request)
		assert.Equal(t, "", provider(req))
	})
}

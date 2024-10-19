package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsGetter(t *testing.T) {
	generator := func() string { return "generated" }

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		options := NewOptions(generator, "X-Request-ID", "")
		getter := options.Getter()
		req := &http.Request{Header: http.Header{}}
		req.Header.Set("X-Request-ID", "in-header")
		assert.Equal(t, "in-header", getter(req))
	})
	t.Run("request without X-Request-ID header", func(t *testing.T) {
		options := NewOptions(generator, "", "")
		getter := options.Getter()
		assert.Equal(t, "generated", getter(new(http.Request)))
	})
}

func TestOptionsResponseSetter(t *testing.T) {
	t.Run("response with X-Request-ID header", func(t *testing.T) {
		options := NewOptions(nil, "", "X-Request-ID")
		respSetter := options.ResponseSetter()
		w := httptest.NewRecorder()
		respSetter(w, "test1")
		assert.Equal(t, "test1", w.Header().Get("X-Request-ID"))
	})
	t.Run("response without X-Request-ID header", func(t *testing.T) {
		options := NewOptions(nil, "", "")
		respSetter := options.ResponseSetter()
		w := httptest.NewRecorder()
		respSetter(w, "test2")
		assert.Empty(t, w.Header().Get("X-Request-ID"))
	})
}

func TestOptionsHandler(t *testing.T) {
	generator := func() string { return "generated" }
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!\n"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		options := NewOptions(generator, "X-Request-ID", "X-Request-ID")
		mockHandler := options.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "in-header", r.Header.Get("X-Request-ID"))
			assert.Equal(t, "in-header", Get(r.Context()))
			baseHandler.ServeHTTP(w, r)
		}))
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Request-ID", "in-header")

		w := httptest.NewRecorder()
		mockHandler.ServeHTTP(w, req)

		assert.Equal(t, "in-header", w.Header().Get("X-Request-ID"))
	})
	t.Run("request without X-Request-ID header", func(t *testing.T) {
		options := NewOptions(generator, "", "X-Request-ID")
		mockHandler := options.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "", r.Header.Get("X-Request-ID"))
			assert.Equal(t, "generated", Get(r.Context()))
			baseHandler.ServeHTTP(w, r)
		}))
		req := httptest.NewRequest("GET", "/", nil)

		w := httptest.NewRecorder()
		mockHandler.ServeHTTP(w, req)

		assert.Equal(t, "generated", w.Header().Get("X-Request-ID"))
	})
}
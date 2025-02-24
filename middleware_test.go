package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestOptions(generator generator, requestHeader string, responseHeader string) *Options {
	return &Options{
		generator:      generator,
		requestHeader:  requestHeader,
		responseHeader: responseHeader,
	}
}

func TestMiddlewareGetter(t *testing.T) {
	generator := func() string { return "generated" }

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		mw := newMiddleware(newTestOptions(generator, "X-Request-ID", ""))
		req := &http.Request{Header: http.Header{}}
		req.Header.Set("X-Request-ID", "in-header")
		assert.Equal(t, "in-header", mw.header.provider(req))
	})
	t.Run("request without X-Request-ID header", func(t *testing.T) {
		mw := newMiddleware(newTestOptions(generator, "", ""))
		assert.Equal(t, "generated", mw.header.provider(new(http.Request)))
	})
}

func TestMiddlewareResponseSetter(t *testing.T) {
	t.Run("response with X-Request-ID header", func(t *testing.T) {
		mw := newMiddleware(newTestOptions(nil, "", "X-Request-ID"))
		respSetter := mw.header.responseSetter
		w := httptest.NewRecorder()
		respSetter(w, "test1")
		assert.Equal(t, "test1", w.Header().Get("X-Request-ID"))
	})
	t.Run("response without X-Request-ID header", func(t *testing.T) {
		mw := newMiddleware(newTestOptions(nil, "", ""))
		respSetter := mw.header.responseSetter
		w := httptest.NewRecorder()
		respSetter(w, "test2")
		assert.Empty(t, w.Header().Get("X-Request-ID"))
	})
}

func TestMiddlewareHandler(t *testing.T) {
	generator := func() string { return "generated" }
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Hello, world!\n"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		mw := newMiddleware(newTestOptions(generator, "X-Request-ID", "X-Request-ID"))
		mockHandler := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "in-header", r.Header.Get("X-Request-ID"))
			assert.Equal(t, "in-header", mw.header.Get(r.Context()))
			baseHandler.ServeHTTP(w, r)
		}))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", "in-header")

		w := httptest.NewRecorder()
		mockHandler.ServeHTTP(w, req)

		assert.Equal(t, "in-header", w.Header().Get("X-Request-ID"))
	})
	t.Run("request without X-Request-ID header", func(t *testing.T) {
		mw := newMiddleware(newTestOptions(generator, "", "X-Request-ID"))
		mockHandler := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "", r.Header.Get("X-Request-ID"))
			assert.Equal(t, "generated", mw.header.Get(r.Context()))
			baseHandler.ServeHTTP(w, r)
		}))
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		w := httptest.NewRecorder()
		mockHandler.ServeHTTP(w, req)

		assert.Equal(t, "generated", w.Header().Get("X-Request-ID"))
	})
}

func TestMiddlewareGetNamespace(t *testing.T) {
	mw := newMiddleware(newTestOptions(nil, "", ""))
	assert.NotNil(t, mw.GetNamespace())
}

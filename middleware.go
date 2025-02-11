package requestid

import (
	"log/slog"
	"net/http"

	"github.com/akm/slogctx"
)

// Middleware provides a middleware for generating and setting request ID.
type Middleware struct {
	namespace      *slogctx.Namespace
	provider       provider
	responseSetter responseSetter
}

// New returns a new Namespace from the given options.
func New(opts ...Option) *Middleware {
	options := newDefaultOptions()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newMiddleware(options)
}

func newMiddleware(options *Options) *Middleware {
	return &Middleware{
		namespace:      getSlogctxNamespace(options.slogctxNamespace, options.logAttr),
		provider:       newProvider(options.generator, options.requestHeader),
		responseSetter: newResponseSetter(options.responseHeader),
	}
}

// Wrap wraps the given http.Handler with the middleware.
// The middleware generates a request ID and sets it to the request context.
func (f *Middleware) Wrap(h http.Handler) http.Handler {
	return wrapHTTPHandler(h, f.provider, f.responseSetter)
}

// NewLogger returns a new logger with the Middleware.
func (f *Middleware) NewLogger(h slog.Handler) *slog.Logger {
	return f.namespace.New(h)
}

// GetNamespace returns the slogctx.Namespace for logging with request ID.
func (f *Middleware) GetNamespace() *slogctx.Namespace {
	return f.namespace
}

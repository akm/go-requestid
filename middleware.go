package requestid

import (
	"log/slog"
	"net/http"

	"github.com/akm/slogctx"
)

// Middleware provides a middleware for generating and setting request ID.
type Middleware struct {
	namespace *slogctx.Namespace
	header    *Header
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
	header := newHeader(&HeaderOptions{
		logAttr:        options.logAttr,
		provider:       newProvider(options.generator, options.requestHeader),
		responseSetter: newResponseSetter(options.responseHeader),
	})
	ctxNs := options.slogctxNamespace
	if ctxNs == nil {
		ctxNs = slogctx.NewNamespace()
	}
	header.addRecordConvTo(ctxNs)

	return &Middleware{
		namespace: ctxNs,
		header:    header,
	}
}

// Wrap wraps the given http.Handler with the middleware.
// The middleware generates a request ID and sets it to the request context.
func (f *Middleware) Wrap(h http.Handler) http.Handler {
	return f.header.Wrap(h)
}

// NewLogger returns a new logger with the Middleware.
func (f *Middleware) NewLogger(h slog.Handler) *slog.Logger {
	return f.namespace.New(h)
}

// WrapSlogHandler wraps the given slog.Handler with the Middleware.
func (f *Middleware) WrapSlogHandler(h slog.Handler) slog.Handler {
	return f.namespace.Wrap(h)
}

// GetNamespace returns the slogctx.Namespace for logging with request ID.
func (f *Middleware) GetNamespace() *slogctx.Namespace {
	return f.namespace
}

package requestid

import (
	"net/http"

	"github.com/akm/slogctx"
)

// Namespace provides a middleware for generating and setting request ID.
type Namespace struct {
	SlogctxNamespace *slogctx.Namespace
	provider         provider
	responseSetter   func(w http.ResponseWriter, id string)
}

// New returns a new Namespace from the given options.
func New(opts ...Option) *Namespace {
	options := newDefaultOptions()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newNamespace(options)
}

func newNamespace(options *Options) *Namespace {
	return &Namespace{
		SlogctxNamespace: getSlogctxNamespace(options.slogctxNamespace, options.logAttr),
		provider:         newProvider(options.generator, options.requestHeader),
		responseSetter:   newResponseSetter(options.responseHeader),
	}
}

// Wrap wraps the given http.Handler with the middleware.
// The middleware generates a request ID and sets it to the request context.
func (f *Namespace) Wrap(h http.Handler) http.Handler {
	return wrapHttpHandler(h, f.provider, f.responseSetter)
}

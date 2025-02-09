package requestid

import (
	"net/http"

	"github.com/akm/slogctx"
)

type Namespace struct {
	SlogctxNamespace *slogctx.Namespace
	provider         provider
	responseSetter   func(w http.ResponseWriter, id string)
}

func New(opts ...Option) *Namespace {
	options := newDefaultOptions()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newNamespace(options)
}

func newNamespace(options *Options) *Namespace {
	var slogctxNamespace *slogctx.Namespace
	if options.slogctxNamespace != nil {
		slogctxNamespace = options.slogctxNamespace
	} else {
		slogctxNamespace = slogctx.NewNamespace()
	}
	slogctxNamespace.AddRecordConv(RecordConv(options.logAttr))
	return &Namespace{
		SlogctxNamespace: slogctxNamespace,
		provider:         newProvider(options.generator, options.requestHeader),
		responseSetter:   newResponseSetter(options.responseHeader),
	}
}

func newProvider(generator generator, requestHeader string) provider {
	coreProvider := generatorProvider(generator)
	if requestHeader != "" {
		return requestIdProviderWrapper(coreProvider, requestHeader)
	} else {
		return coreProvider
	}
}

func newResponseSetter(responseHeader string) func(w http.ResponseWriter, id string) {
	if responseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(responseHeader, id)
		}
	} else {
		return func(http.ResponseWriter, string) {}
	}
}

func (f *Namespace) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := f.provider(r)
		ctx := newContext(r.Context(), requestID)
		f.responseSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

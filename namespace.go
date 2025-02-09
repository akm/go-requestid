package requestid

import (
	"net/http"

	"github.com/akm/slogctx"
)

type Namespace struct {
	options        *Options
	SlogwNamespace *slogctx.Namespace
}

func New(opts ...Option) *Namespace {
	options := newDefaultOptions()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newFactory(options)
}

func newFactory(options *Options) *Namespace {
	var slogwNamespace *slogctx.Namespace
	if options.slogctxNamespace != nil {
		slogwNamespace = options.slogctxNamespace
	} else {
		slogwNamespace = slogctx.NewNamespace()
	}
	return &Namespace{
		options:        options,
		SlogwNamespace: slogwNamespace,
	}
}

func (f *Namespace) getter() provider {
	coreProvider := generatorProvider(f.options.generator)
	if f.options.requestHeader != "" {
		return requestIdProviderWrapper(coreProvider, f.options.requestHeader)
	} else {
		return coreProvider
	}
}

func (f *Namespace) responseSetter() func(w http.ResponseWriter, id string) {
	if f.options.responseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(f.options.responseHeader, id)
		}
	} else {
		return func(http.ResponseWriter, string) {}
	}
}

func (f *Namespace) Wrap(h http.Handler) http.Handler {
	getter := f.getter()
	respSetter := f.responseSetter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := getter(r)
		ctx := set(r.Context(), requestID)
		respSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

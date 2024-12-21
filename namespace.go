package requestid

import (
	"net/http"
)

type Namespace struct {
	options *Options
}

func newFactory(options *Options) *Namespace {
	return &Namespace{options: options}
}

func (f *Namespace) getter() provider {
	coreProvider := generatorProvider(f.options.Generator)
	if f.options.RequestHeader != "" {
		return requestIdProviderWrapper(coreProvider, f.options.RequestHeader)
	} else {
		return coreProvider
	}
}

func (f *Namespace) responseSetter() func(w http.ResponseWriter, id string) {
	if f.options.ResponseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(f.options.ResponseHeader, id)
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

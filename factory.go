package requestid

import (
	"net/http"
)

type Factory struct {
	generator      generator
	requestHeader  string
	responseHeader string
}

func (f *Factory) getter() Provider {
	coreProvider := GeneratorProvider(f.generator)
	if f.requestHeader != "" {
		return RequestIdProviderWrapper(coreProvider, f.requestHeader)
	} else {
		return coreProvider
	}
}

func (f *Factory) responseSetter() func(w http.ResponseWriter, id string) {
	if f.responseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(f.responseHeader, id)
		}
	} else {
		return func(http.ResponseWriter, string) {}
	}
}

func (f *Factory) Wrap(h http.Handler) http.Handler {
	getter := f.getter()
	respSetter := f.responseSetter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := getter(r)
		ctx := set(r.Context(), requestID)
		respSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

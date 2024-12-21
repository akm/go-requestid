package requestid

import (
	"net/http"
)

type Factory struct {
	options *Options
}

func newFactory(options *Options) *Factory {
	return &Factory{options: options}
}

func (f *Factory) getter() Provider {
	coreProvider := GeneratorProvider(f.options.Generator)
	if f.options.RequestHeader != "" {
		return RequestIdProviderWrapper(coreProvider, f.options.RequestHeader)
	} else {
		return coreProvider
	}
}

func (f *Factory) responseSetter() func(w http.ResponseWriter, id string) {
	if f.options.ResponseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(f.options.ResponseHeader, id)
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

package requestid

import (
	"net/http"
)

type Factory struct {
	generator      Generator
	requestHeader  string
	responseHeader string
}

func NewFactory(generator Generator, requestHeader string, responseHeader string) *Factory {
	return &Factory{
		generator:      generator,
		requestHeader:  requestHeader,
		responseHeader: responseHeader,
	}
}

func (f *Factory) Getter() Provider {
	coreProvider := GeneratorProvider(f.generator)
	if f.requestHeader != "" {
		return RequestIdProviderWrapper(coreProvider, f.requestHeader)
	} else {
		return coreProvider
	}
}

func (f *Factory) ResponseSetter() func(w http.ResponseWriter, id string) {
	if f.responseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(f.responseHeader, id)
		}
	} else {
		return func(http.ResponseWriter, string) {}
	}
}

func (f *Factory) Handler(h http.Handler) http.Handler {
	getter := f.Getter()
	respSetter := f.ResponseSetter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := getter(r)
		ctx := set(r.Context(), requestID)
		respSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

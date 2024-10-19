package requestid

import (
	"net/http"
)

type Options struct {
	generator      Generator
	requestHeader  string
	responseHeader string
}

func NewOptions(generator Generator, requestHeader string, responseHeader string) *Options {
	return &Options{
		generator:      generator,
		requestHeader:  requestHeader,
		responseHeader: responseHeader,
	}
}

func (f *Options) Getter() Provider {
	coreProvider := GeneratorProvider(f.generator)
	if f.requestHeader != "" {
		return RequestIdProviderWrapper(coreProvider, f.requestHeader)
	} else {
		return coreProvider
	}
}

func (f *Options) ResponseSetter() func(w http.ResponseWriter, id string) {
	if f.responseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(f.responseHeader, id)
		}
	} else {
		return func(http.ResponseWriter, string) {}
	}
}

func (f *Options) Handler(h http.Handler) http.Handler {
	getter := f.Getter()
	respSetter := f.ResponseSetter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := getter(r)
		ctx := set(r.Context(), requestID)
		respSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

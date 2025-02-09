package requestid

import "net/http"

func wrapHttpHandler(h http.Handler, provider provider, responseSetter func(w http.ResponseWriter, id string)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := provider(r)
		ctx := newContext(r.Context(), requestID)
		responseSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type provider = func(req *http.Request) string

func newProvider(generator generator, requestHeader string) provider {
	coreProvider := generatorProvider(generator)
	if requestHeader != "" {
		return requestIdProviderWrapper(coreProvider, requestHeader)
	} else {
		return coreProvider
	}
}

func generatorProvider(generator generator) provider {
	return func(_ *http.Request) string { return generator() }
}

func requestIdProviderWrapper(next provider, requestHeader string) provider {
	return func(req *http.Request) string {
		if requestID := req.Header.Get(requestHeader); requestID != "" {
			return requestID
		}
		return next(req)
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

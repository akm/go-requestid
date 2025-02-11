package requestid

import "net/http"

type processor struct {
	provider       provider
	responseSetter responseSetter
}

func newProcessor(provider provider, responseSetter responseSetter) *processor {
	return &processor{
		provider:       provider,
		responseSetter: responseSetter,
	}
}

func (p *processor) wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := p.provider(r)
		ctx := newContext(r.Context(), requestID)
		p.responseSetter(w, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type provider = func(req *http.Request) string

func newProvider(generator generator, requestHeader string) provider {
	coreProvider := generatorProvider(generator)
	if requestHeader != "" {
		return requestIDProviderWrapper(coreProvider, requestHeader)
	}
	return coreProvider
}

func generatorProvider(generator generator) provider {
	return func(_ *http.Request) string { return generator() }
}

func requestIDProviderWrapper(next provider, requestHeader string) provider {
	return func(req *http.Request) string {
		if requestID := req.Header.Get(requestHeader); requestID != "" {
			return requestID
		}
		return next(req)
	}
}

type responseSetter = func(w http.ResponseWriter, id string)

func newResponseSetter(responseHeader string) responseSetter {
	if responseHeader != "" {
		return func(w http.ResponseWriter, id string) {
			w.Header().Set(responseHeader, id)
		}
	}
	return func(http.ResponseWriter, string) {}
}

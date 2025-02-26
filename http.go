package requestid

import "net/http"

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

func requestIDProvider(requestHeader string) provider {
	return func(req *http.Request) string {
		if requestID := req.Header.Get(requestHeader); requestID != "" {
			return requestID
		}
		return ""
	}
}

func selectNotEmpty(providers ...provider) provider {
	return func(req *http.Request) string {
		for _, p := range providers {
			if id := p(req); id != "" {
				return id
			}
		}
		return ""
	}
}

func requestIDProviderWrapper(next provider, requestHeader string) provider {
	return selectNotEmpty(requestIDProvider(requestHeader), next)
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

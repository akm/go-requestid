package requestid

import "net/http"

type Provider = func(req *http.Request) string

func RequestIdProviderWrapper(next Provider, requestHeader string) Provider {
	return func(req *http.Request) string {
		if requestID := req.Header.Get(requestHeader); requestID != "" {
			return requestID
		}
		return next(req)
	}
}

func GeneratorProvider(generator Generator) Provider {
	return func(_ *http.Request) string { return generator() }
}

package requestid

import "net/http"

type provider = func(req *http.Request) string

func requestIdProviderWrapper(next provider, requestHeader string) provider {
	return func(req *http.Request) string {
		if requestID := req.Header.Get(requestHeader); requestID != "" {
			return requestID
		}
		return next(req)
	}
}

func generatorProvider(generator generator) provider {
	return func(_ *http.Request) string { return generator() }
}

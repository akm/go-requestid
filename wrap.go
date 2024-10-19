package requestid

import (
	"net/http"
)

func Wrap(next http.Handler, opts ...Option) http.Handler {
	options := Default()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newFactory(options).Wrap(next)
}

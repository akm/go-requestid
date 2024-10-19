package requestid

import (
	"net/http"
)

func Wrap(next http.Handler) http.Handler {
	headerName := "X-Request-ID"
	factory := newOptions(defaultGenerator, headerName, headerName)
	return factory.Wrap(next)
}

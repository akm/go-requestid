package requestid

import (
	"net/http"
)

func Wrap(next http.Handler) http.Handler {
	headerName := "X-Request-ID"
	factory := NewOptions(defaultGenerator, headerName, headerName)
	return factory.Handler(next)
}

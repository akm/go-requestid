package requestid

import (
	"net/http"
)

func Handler(next http.Handler) http.Handler {
	headerName := "X-Request-ID"
	factory := NewFactory(defaultGenerator, headerName, headerName)
	return factory.Handler(next)
}

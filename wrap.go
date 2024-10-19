package requestid

import (
	"net/http"
)

func Wrap(next http.Handler) http.Handler {
	return Default().Wrap(next)
}

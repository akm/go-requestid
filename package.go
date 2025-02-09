package requestid

import (
	"net/http"
)

var defaultNamespace = newNamespace(newDefaultOptions())

func Default() *Namespace {
	return defaultNamespace
}

func SetDefault(ns *Namespace) {
	defaultNamespace = ns
}

func Wrap(next http.Handler) http.Handler {
	return defaultNamespace.Wrap(next)
}

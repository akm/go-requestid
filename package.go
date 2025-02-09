package requestid

import (
	"net/http"
)

var defaultNamespace = newNamespace(newDefaultOptions())

func DefaultNamespace() *Namespace {
	return defaultNamespace
}

func SetDefaultNamespace(ns *Namespace) {
	defaultNamespace = ns
}

func Wrap(next http.Handler, opts ...Option) http.Handler {
	var ns *Namespace
	if len(opts) == 0 {
		ns = defaultNamespace
	} else {
		ns = New(opts...)
	}
	return ns.Wrap(next)
}

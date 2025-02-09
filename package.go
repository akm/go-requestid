package requestid

import (
	"net/http"
)

var degaultNamespace = newFactory(newDefaultOptions())

func DefaultNamespace() *Namespace {
	return degaultNamespace
}

func SetDefaultNamespace(ns *Namespace) {
	degaultNamespace = ns
}

func Wrap(next http.Handler, opts ...Option) http.Handler {
	var ns *Namespace
	if len(opts) == 0 {
		ns = degaultNamespace
	} else {
		ns = New(opts...)
	}
	return ns.Wrap(next)
}

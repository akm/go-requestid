package requestid

import (
	"net/http"
)

var degaultNamespace = newFactory(Default())

func DefaultNamespace() *Namespace {
	return degaultNamespace
}

func SetDefaultNamespace(ns *Namespace) {
	degaultNamespace = ns
}

func RegisterSlogHandle(key string) {
	degaultNamespace.RegisterSlogHandle(key)
}

func Wrap(next http.Handler, opts ...Option) http.Handler {
	options := Default()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newFactory(options).Wrap(next)
}

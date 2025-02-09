package requestid

import (
	"log/slog"
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

func NewLogger(h slog.Handler) *slog.Logger {
	return defaultNamespace.SlogctxNamespace.New(h)
}

package requestid

import (
	"log/slog"
	"net/http"
)

var defaultNamespace = newNamespace(newDefaultOptions())

// Default returns the default Namespace.
func Default() *Namespace {
	return defaultNamespace
}

// SetDefault sets the default Namespace.
func SetDefault(ns *Namespace) {
	defaultNamespace = ns
}

// Wrap wraps the given http.Handler with the default Namespace.
func Wrap(next http.Handler) http.Handler {
	return defaultNamespace.Wrap(next)
}

// NewLogger returns a new logger with the default Namespace.
func NewLogger(h slog.Handler) *slog.Logger {
	return defaultNamespace.SlogctxNamespace.New(h)
}

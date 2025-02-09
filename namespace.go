package requestid

import (
	"net/http"

	"github.com/akm/slogctx"
)

type Namespace struct {
	SlogctxNamespace *slogctx.Namespace
	provider         provider
	responseSetter   func(w http.ResponseWriter, id string)
}

func New(opts ...Option) *Namespace {
	options := newDefaultOptions()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newNamespace(options)
}

func newNamespace(options *Options) *Namespace {
	var slogctxNamespace *slogctx.Namespace
	if options.slogctxNamespace != nil {
		slogctxNamespace = options.slogctxNamespace
	} else {
		slogctxNamespace = slogctx.NewNamespace()
	}
	slogctxNamespace.AddRecordConv(RecordConv(options.logAttr))
	return &Namespace{
		SlogctxNamespace: slogctxNamespace,
		provider:         newProvider(options.generator, options.requestHeader),
		responseSetter:   newResponseSetter(options.responseHeader),
	}
}

func (f *Namespace) Wrap(h http.Handler) http.Handler {
	return wrapHttpHandler(h, f.provider, f.responseSetter)
}

package requestid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	h := newHeader(newHeaderOptions())
	assert.Equal(t, "", h.Get(ctx))
	ctx = h.newContext(ctx, "test")
	assert.Equal(t, "test", h.Get(ctx))
}

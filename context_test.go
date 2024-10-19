package requestid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "", Get(ctx))
	ctx = set(ctx, "test")
	assert.Equal(t, "test", Get(ctx))
}

package requestid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenerator(t *testing.T) {
	generator := defaultGenerator
	require.NotNil(t, generator)

	sampleCount := 100
	sample := make([]string, sampleCount)
	for i := 0; i < sampleCount; i++ {
		sample[i] = generator()
	}

	t.Run("length", func(t *testing.T) {
		for _, s := range sample {
			assert.Len(t, s, 8)
		}
	})
	t.Run("unique", func(t *testing.T) {
		unique := make(map[string]struct{})
		for _, s := range sample {
			unique[s] = struct{}{}
		}
		assert.Len(t, unique, sampleCount)
	})
}

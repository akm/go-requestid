package requestid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	assert.Equal(t, degaultNamespace, DefaultNamespace())
	assert.Equal(t, degaultNamespace, DefaultNamespace(), "DefaultNamespace should return the same instance")

	ns1 := New(ResponseHeader("X-REQ-ID"))
	SetDefaultNamespace(ns1)
	defer SetDefaultNamespace(degaultNamespace)

	assert.Equal(t, ns1, DefaultNamespace())
	assert.Equal(t, ns1, DefaultNamespace(), "DefaultNamespace should return the same instance")
}

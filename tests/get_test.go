package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	out, _, err := epinio.Run("get")
	assert.NoError(t, err)

	assert.Contains(t, out, "Usage:")
}

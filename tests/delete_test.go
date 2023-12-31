package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	out, _, err := epinio.Run("delete")
	assert.NoError(t, err)

	assert.Contains(t, out, "Usage:")
}

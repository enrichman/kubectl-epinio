package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	_, stderr, err := epinio.Create("user", "foo")
	require.NoError(t, err, stderr)

	// check user exists
	stdout, _, err := epinio.Get("user", "foo")
	require.NoError(t, err)
	outTable := parseOutTable(stdout)
	assert.Equal(t, "foo", outTable[1][0])

	_, _, err = epinio.Delete("user", "foo", "-y")
	assert.NoError(t, err, stderr)

	// check user do not exists
	stdout, _, err = epinio.Get("user", "foo")
	require.NoError(t, err)
	outTable = parseOutTable(stdout)
	assert.Equal(t, 1, len(outTable))
}

func TestDeleteNonExistingUser(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	// check user do not exists
	stdout, _, err := epinio.Get("user", "foo")
	require.NoError(t, err)
	outTable := parseOutTable(stdout)
	assert.Equal(t, 1, len(outTable), stdout)

	_, _, err = epinio.Delete("user", "foo")
	assert.Error(t, err)
}

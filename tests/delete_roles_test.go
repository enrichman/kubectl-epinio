package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteRole(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	_, _, err = epinio.Create("role", "foo")
	assert.NoError(t, err)

	// check role exists
	stdout, _, err := epinio.Get("role", "foo")
	require.NoError(t, err)
	outTable := parseOutTable(stdout)
	assert.Equal(t, "foo", outTable[1][0])

	_, _, err = epinio.Delete("role", "foo", "-y")
	require.NoError(t, err)

	// check role do not exists
	stdout, _, err = epinio.Get("role", "foo")
	require.NoError(t, err)
	outTable = parseOutTable(stdout)
	assert.Equal(t, 1, len(outTable), stdout)
}

func TestDeleteNonExistingRole(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	// check role do not exists
	stdout, _, err := epinio.Get("role", "foo")
	require.NoError(t, err)
	outTable := parseOutTable(stdout)
	assert.Equal(t, 1, len(outTable), outTable)

	_, _, err = epinio.Delete("role", "foo")
	assert.Error(t, err)
}

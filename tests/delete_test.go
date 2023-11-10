package tests

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteRole(t *testing.T) {
	testCases := []struct {
		name            string
		createdRoleName string
		args            []string
		expectErr       bool
	}{
		{
			name:            "delete existing role",
			createdRoleName: "foo",
			args:            []string{"foo"},
			expectErr:       false,
		},
		{
			name:            "delete nonexistent role",
			createdRoleName: "foo",
			args:            []string{"bar"},
			expectErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createArgs := []string{"create", "role", tc.createdRoleName}
			cmd := exec.Command(cmdExecPath(t), createArgs...)

			_, err := cmd.Output()
			require.NoError(t, err)

			if len(tc.args) == 0 || tc.createdRoleName != tc.args[0] {
				defer deleteRole(t, []string{tc.createdRoleName}, false) // clean up created role
			}

			checkRoleExists(t, tc.createdRoleName, true)

			deleteRole(t, tc.args, tc.expectErr)

			checkRoleExists(t, tc.args[0], false)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name            string
		createdUserName string
		args            []string
		expectErr       bool
	}{
		{
			name:            "delete existing user",
			createdUserName: "foo",
			args:            []string{"foo"},
			expectErr:       false,
		},
		{
			name:            "delete nonexistent user",
			createdUserName: "foo",
			args:            []string{"bar"},
			expectErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createArgs := []string{"create", "user", tc.createdUserName}
			cmd := exec.Command(cmdExecPath(t), createArgs...)

			_, err := cmd.Output()
			require.NoError(t, err)

			if len(tc.args) == 0 || tc.createdUserName != tc.args[0] {
				defer deleteUser(t, []string{tc.createdUserName}, false) // clean up created user
			}

			checkUserExists(t, tc.createdUserName, true)

			deleteUser(t, tc.args, tc.expectErr)

			checkUserExists(t, tc.args[0], false)
		})
	}
}

func deleteRole(t *testing.T, args []string, expectError bool) {
	deleteEntity(t, "role", args, expectError)
}

func deleteUser(t *testing.T, args []string, expectError bool) {
	deleteEntity(t, "user", args, expectError)
}

func deleteEntity(t *testing.T, entityName string, args []string, expectError bool) {
	entityArg := validateEntity(t, entityName)

	deleteArgs := []string{"delete", entityArg}
	deleteArgs = append(deleteArgs, args...)
	deleteArgs = append(deleteArgs, "-y") // non-interactive mode

	cmd := exec.Command(cmdExecPath(t), deleteArgs...)

	_, err := cmd.Output()
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func checkRoleExists(t *testing.T, name string, shouldBeFound bool) {
	checkEntityExists(t, "role", name, shouldBeFound)
}

func checkUserExists(t *testing.T, name string, shouldBeFound bool) {
	checkEntityExists(t, "user", name, shouldBeFound)
}

func checkEntityExists(t *testing.T, entityName string, name string, shouldBeFound bool) {
	entityArg := validateEntity(t, entityName)

	getArgs := []string{"get", entityArg}
	cmd := exec.Command(cmdExecPath(t), getArgs...)

	out, err := cmd.Output()
	require.NoError(t, err)

	isFound := false

	foundEntries := strings.Split(string(out), "\n")
	for _, entry := range foundEntries {
		if entry == name {
			isFound = true
			break
		}
	}

	assert.Equal(t, shouldBeFound, isFound)
}

func validateEntity(t *testing.T, name string) string {
	switch name {
	case "user", "role":
		return name
	default:
		assert.False(t, true, fmt.Sprintf("Got entity name %s, expecting role or user", name))
	}

	return ""
}

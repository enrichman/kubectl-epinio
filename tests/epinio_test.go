package tests

import (
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cmd := exec.Command(
		cmdExecPath(t),
		"get",
	)

	out, err := cmd.Output()
	assert.NoError(t, err)

	assert.Contains(t, string(out), "Usage:")
}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		expectedEntries [][]interface{}
		expectErr       bool
	}{
		{
			name: "get all users",
			args: []string{},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{"admin", true, 1, ""},
				{"admin@epinio.io", true, 1, ""},
				{"epinio", false, 1, ""},
			},
		},
		{
			name: "get all users",
			args: []string{},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{"admin", true, 1, ""},
				{"admin@epinio.io", true, 1, ""},
				{"epinio", false, 1, ""},
			},
		},
		{
			name: "get single user with exact match",
			args: []string{"epinio"},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{"epinio", false, 1, ""},
			},
		},
		{
			name: "get nonexistent user",
			args: []string{"nothere"},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
			},
		},
		{
			name: "get multiple users with exact matches",
			args: []string{"epinio", "admin"},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{"admin", true, 1, ""},
				{"epinio", false, 1, ""},
			},
		},
		{
			name: "get multiple users with single exact match",
			args: []string{"epinio", "nothere"},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{"epinio", false, 1, ""},
			},
		},
	}

	for _, usersArg := range []string{"user", "users"} {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s/%s", tc.name, usersArg), func(t *testing.T) {
				args := []string{"get", usersArg}
				args = append(args, tc.args...)

				cmd := exec.Command(cmdExecPath(t), args...)

				out, err := cmd.Output()
				assert.NoError(t, err)

				rows := strings.Split(strings.TrimSpace(string(out)), "\n")
				assert.True(t, len(rows) > 0)

				for i, row := range rows {
					rowCells := strings.FieldsFunc(row, func(r rune) bool {
						return r == '\t'
					})
					assert.Equal(t, len(tc.expectedEntries[i]), len(rowCells))

					// check headers
					if i == 0 {
						for j, expected := range tc.expectedEntries[i] {
							assert.Equal(t, expected, rowCells[j])
						}
						continue
					}

					// check values
					for j, expected := range tc.expectedEntries[i] {
						switch j {
						case 0:
							assert.Equal(t, expected, rowCells[j])
						case 1:
							parsedBool, err := strconv.ParseBool(rowCells[j])
							assert.NoError(t, err)
							assert.Equal(t, expected, parsedBool)
						case 2:
							parsedInt, err := strconv.Atoi(rowCells[j])
							assert.NoError(t, err)
							assert.Equal(t, expected, parsedInt)
						case 3:
							_, err := time.ParseDuration(rowCells[j])
							assert.NoError(t, err)
						}
					}
				}
			})
		}
	}
}

func TestGetRoles(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		expectedEntries [][]interface{}
	}{
		{
			name: "get all roles",
			args: []string{},
			expectedEntries: [][]interface{}{
				{"ID", "NAME", "DEFAULT", "ACTIONS", "AGE"},
				{"admin", "Admin Role", false, 0, "-"},
			},
		},
	}

	for _, usersArg := range []string{"role", "roles"} {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s/%s", tc.name, usersArg), func(t *testing.T) {
				args := []string{"get", usersArg}
				args = append(args, tc.args...)

				cmd := exec.Command(cmdExecPath(t), args...)

				out, err := cmd.Output()
				assert.NoError(t, err)

				rows := strings.Split(strings.TrimSpace(string(out)), "\n")
				assert.True(t, len(rows) > 0)

				for i, row := range rows {
					rowCells := strings.FieldsFunc(row, func(r rune) bool {
						return r == '\t'
					})
					assert.Equal(t, len(tc.expectedEntries[i]), len(rowCells))

					// check headers
					if i == 0 {
						for j, expected := range tc.expectedEntries[i] {
							assert.Equal(t, expected, rowCells[j])
						}
						continue
					}

					// check values
					for j, expected := range tc.expectedEntries[i] {
						actual := rowCells[j]

						switch j {
						case 0, 1:
							assert.Equal(t, expected, actual)
						case 2:
							parsedBool, err := strconv.ParseBool(actual)
							assert.NoError(t, err)
							assert.Equal(t, expected, parsedBool)
						case 3:
							parsedInt, err := strconv.Atoi(actual)
							assert.NoError(t, err)
							assert.Equal(t, expected, parsedInt)
						case 4:
							if actual != "-" {
								_, err := time.ParseDuration(actual)
								assert.NoError(t, err)
							}
						}
					}
				}
			})
		}
	}
}

// getGitRepoRoot returns the root of the git repository in which it is executed.
func getGitRepoRoot(t *testing.T) string {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--show-toplevel",
	)
	root, err := cmd.Output()
	assert.NoError(t, err)

	return strings.TrimRight(string(root), "\n")
}

// cmdExecPath returns the path to the kubectl-epinio command binary.
func cmdExecPath(t *testing.T) string {
	return path.Join(getGitRepoRoot(t), "output", "kubectl-epinio")
}

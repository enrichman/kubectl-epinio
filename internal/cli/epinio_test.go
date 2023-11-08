package cli_test

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
	"testing"

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
		expectedEntries []string
		expectErr       bool
	}{
		{
			name: "get all users",
			args: []string{},
			expectedEntries: []string{
				"admin",
				"admin@epinio.io",
				"epinio",
			},
		},
		{
			name: "get all users",
			args: []string{},
			expectedEntries: []string{
				"admin",
				"admin@epinio.io",
				"epinio",
			},
		},
		{
			name: "get single user with exact match",
			args: []string{"epinio"},
			expectedEntries: []string{
				"epinio",
			},
		},
		{
			name:            "get nonexistent user",
			args:            []string{"nothere"},
			expectedEntries: []string{},
		},
		{
			name:            "get multiple users with exact matches",
			args:            []string{"epinio", "admin"},
			expectedEntries: []string{"admin", "epinio"},
		},
		{
			name:            "get multiple users with single exact match",
			args:            []string{"epinio", "nothere"},
			expectedEntries: []string{"epinio"},
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

				entries := strings.Split(string(out), "\n")
				entries = entries[:len(entries)-1] // The last entry is empty

				expectedEntries := []string{"USERNAME"} // header
				expectedEntries = append(expectedEntries, tc.expectedEntries...)

				assert.Equal(t, len(expectedEntries), len(entries))

				for i := range tc.expectedEntries {
					assert.Equal(t, expectedEntries[i], entries[i])
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

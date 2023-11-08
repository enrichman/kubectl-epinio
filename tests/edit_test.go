package tests

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/enrichman/kubectl-epinio/pkg/epinio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// XXX: could we make this a tmp file per test case and make this work with the mock editor?
const tmpFile = "/tmp/edited_user.yaml"

func TestEdit(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		createUser      *epinio.User
		newUserContents string
		expectErr       bool
		expectedOutput  string
	}{
		{
			name:       "edit simple user with valid syntax",
			args:       []string{"foo"},
			createUser: &epinio.User{Username: "foo"},
			newUserContents: `username: foo
password: foo`,
			expectedOutput: "",
		},
		{
			name:            "edit simple user making it empty",
			args:            []string{"foo"},
			createUser:      &epinio.User{Username: "foo"},
			newUserContents: "",
			expectErr:       true,
		},
		{
			name:      "edit nonexistent user",
			args:      []string{"foo"},
			expectErr: true,
		},
		{
			name:       "edit multiple users",
			args:       []string{"foo epinio"},
			createUser: &epinio.User{Username: "foo"},
			expectErr:  true, // only one username can be provided
		},
		{
			name:            "edit user name",
			args:            []string{"foo"},
			createUser:      &epinio.User{Username: "foo"},
			newUserContents: "username: bar",
			expectErr:       false,
		},
	}

	for _, usersArg := range []string{"user", "users"} {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.createUser != nil && len(tc.createUser.Username) > 0 {
					createUser(t, *tc.createUser)

					defer deleteUser(t, []string{tc.createUser.Username}, false) // clean up created user
				}

				// Use a non-interactive mock editor to enable automated testing.
				// This editor uses tmpFile as the source of new file contents.
				err := os.Setenv("EDITOR", path.Join(getGitRepoRoot(t), "tests", "mock_editor.sh"))
				require.NoError(t, err)

				f, err := os.Create(tmpFile) // truncated if it already exists
				require.NoError(t, err)

				_, err = f.Write([]byte(tc.newUserContents))
				require.NoError(t, err)

				if tc.createUser != nil {
					checkPreEditContents(t, *tc.createUser)
				}

				args := []string{"edit", usersArg}
				args = append(args, tc.args...)

				cmd := exec.Command(cmdExecPath(t), args...)

				out, err := cmd.Output()
				if tc.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				if tc.expectedOutput == "" {
					assert.Empty(t, out)
				} else {
					assert.Equal(t, tc.expectedOutput, strings.TrimSpace(string(out)))
				}

				if tc.createUser != nil {
					checkUser(t, *tc.createUser)
				}
			})
		}
	}
}

func createUser(t *testing.T, user epinio.User) {
	createArgs := []string{"create", "user", user.Username}

	if len(user.Namespaces) > 0 {
		ns := []string{"--namespaces", strings.Join(user.Namespaces, ",")}
		createArgs = append(createArgs, ns...)
	}

	if len(user.Roles) > 0 {
		roles := []string{"--roles", strings.Join(user.Roles, ",")}
		createArgs = append(createArgs, roles...)
	}

	if len(user.Password) > 0 {
		createArgs = append(createArgs, []string{"--password", user.Password}...)
	}

	cmd := exec.Command(cmdExecPath(t), createArgs...)

	_, err := cmd.Output()
	require.NoError(t, err)
}

// checkUser runs `kubectl-epinio describe user <user_name>` and checks that the output matches user.
func checkUser(t *testing.T, user epinio.User) {
	describeArgs := []string{"describe", "user", user.Username}
	cmd := exec.Command(cmdExecPath(t), describeArgs...)

	out, err := cmd.Output()
	require.NoError(t, err)

	format := "%-15s %s\n"
	var expected strings.Builder
	fmt.Fprintf(&expected, format, "Username:", user.Username)
	fmt.Fprintf(&expected, format, "Password:", user.Password)

	if len(user.Roles) == 0 {
		fmt.Fprintf(&expected, format, "Roles:", "")
		expected.WriteString("\n")
	} else {
		for _, role := range user.Roles {
			fmt.Fprintf(&expected, format, "", role)
		}
	}

	if len(user.Namespaces) == 0 {
		fmt.Fprintf(&expected, format, "Namespaces:", "")
		expected.WriteString("\n")
	} else {
		for _, ns := range user.Namespaces {
			fmt.Fprintf(&expected, format, "", ns)
		}
	}

	assert.Equal(t, expected.String(), string(out))
}

// checkPreEditContents reads the preEditFile and validates its contents for the provided user.
func checkPreEditContents(t *testing.T, user epinio.User) {
	var expected strings.Builder

	fmt.Fprintf(&expected, "username : %s\n", user.Username)
	fmt.Fprintf(&expected, "password:%s\n", user.Password)

	if user.Role != "" {
		fmt.Fprintf(&expected, "role: %s\n", user.Role)
	}

	if len(user.Roles) > 0 {
		fmt.Fprint(&expected, "roles:\n")
	}
	for _, role := range user.Roles {
		fmt.Fprintf(&expected, "  - %s\n", role)
	}

	if len(user.Namespaces) > 0 {
		fmt.Fprint(&expected, "namespaces:\n")
	}
	for _, ns := range user.Namespaces {
		fmt.Fprintf(&expected, "  - %s\n", ns)
	}

	file, err := os.Open(preEditFile)
	require.NoError(t, err)
	defer file.Close()

	commentRegex := regexp.MustCompile(`^#.*$`)

	var preEdit strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		if commentRegex.MatchString(row) {
			// Skip comments, we are only interested in user specs
			continue
		}

		preEdit.WriteString(row + "\n")
	}

	assert.Equal(t, expected.String(), preEdit.String())
}

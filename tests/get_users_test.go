package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

	admin1 := createUser(t, epinio, "admin01", []string{"admin"})
	admin2 := createUser(t, epinio, "admin02", []string{"admin"})
	user1 := createUser(t, epinio, "user01", []string{"user"})
	user2 := createUser(t, epinio, "user02", []string{"admin:workspace", "user"})
	user3 := createUser(t, epinio, "user03", nil)

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
				{admin1, true, 1, ""},
				{admin2, true, 1, ""},
				{user1, false, 1, ""},
				{user2, false, 2, ""},
				{user3, false, 0, ""},
			},
		},
		{
			name: "get single user with exact match",
			args: []string{user1},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{user1, false, 1, ""},
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
			args: []string{user1, admin1},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{user1, false, 1, ""},
				{admin1, true, 1, ""},
			},
		},
		{
			name: "get multiple users with single exact match",
			args: []string{user1, "nothere"},
			expectedEntries: [][]interface{}{
				{"USERNAME", "ADMIN", "ROLES", "AGE"},
				{user1, false, 1, ""},
			},
		},
	}

	for _, resource := range []string{"user", "users"} {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s/%s", tc.name, resource), func(t *testing.T) {
				stdout, _, err := epinio.Get(resource, tc.args...)
				assert.NoError(t, err)

				outTable := parseOutTable(stdout)
				assert.Equal(t, len(tc.expectedEntries), len(outTable), outTable)

				for i, rowCells := range outTable {
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
						case 0:
							assert.Equal(t, expected, actual)
						case 1:
							parsedBool, err := strconv.ParseBool(actual)
							assert.NoError(t, err)
							assert.Equal(t, expected, parsedBool)
						case 2:
							parsedInt, err := strconv.Atoi(actual)
							assert.NoError(t, err)
							assert.Equal(t, expected, parsedInt)
						case 3:
							_, err := time.ParseDuration(actual)
							assert.NoError(t, err)
						}
					}
				}

				if t.Failed() {
					t.Log("Test failed.\nOutput was:\n", stdout)
				}
			})
		}
	}
}

func createUser(t *testing.T, epinio *KubectlEpinio, namePrefix string, roles []string) string {
	t.Helper()

	name := fmt.Sprintf("epinio-user-%s-%s", namePrefix, RandStringBytes(5))
	args := []string{name}

	for _, r := range roles {
		args = append(args, "--roles", r)
	}
	_, _, err := epinio.Create("user", args...)
	require.NoError(t, err)

	return name
}

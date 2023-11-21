package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

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
				{"epinio", false, 1, ""},
				{"admin", true, 1, ""},
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

	for _, resource := range []string{"user", "users"} {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s/%s", tc.name, resource), func(t *testing.T) {
				stdout, _, err := epinio.Get(resource, tc.args...)
				assert.NoError(t, err)

				outTable := parseOutTable(stdout)
				assert.True(t, len(outTable) > 0)

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

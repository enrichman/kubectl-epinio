package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRoles(t *testing.T) {
	epinio, err := NewKubectlEpinio()
	assert.NoError(t, err)

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

	for _, resource := range []string{"role", "roles"} {
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

				if t.Failed() {
					t.Log("Test failed.\nOutput was:\n", stdout)
				}
			})
		}
	}
}

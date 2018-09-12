package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InvalidJSON(t *testing.T) {
	_, err := ParseFilters("{{{}")
	assert.Error(t, err)
}

func Test_Filter(t *testing.T) {
	emptyString := ""
	emptyObject := "{}"

	databaseAllCollections := "{\"db1\":null}"
	databaseNoCollections := "{\"db1\":[]}"
	databaseNamedCollections := "{\"db1\":[\"col1\",\"col3\"]}"

	testCases := []struct {
		filterJSON string
		dbName     string
		collName   string
		expBool    bool
	}{
		{
			// empty filter string should always return true
			emptyString,
			"databaseName",
			"collectionName",
			true,
		},
		{
			// empty json object filter should always return false
			emptyObject,
			"databaseName",
			"collectionName",
			false,
		},

		{
			// special case for checking ONLY database name
			databaseNoCollections,
			"db1",
			"",
			true,
		},

		{
			databaseAllCollections,
			"database1",
			"collectionName",
			false,
		},
		{
			databaseAllCollections,
			"db1",
			"collectionName",
			true,
		},
		{
			databaseNoCollections,
			"db1",
			"collectionName",
			false,
		},
		{
			databaseNamedCollections,
			"db1",
			"collectionName",
			false,
		},
		{
			databaseNamedCollections,
			"db1",
			"col3",
			true,
		},
	}

	for _, tc := range testCases {
		dbFilter, _ := ParseFilters(tc.filterJSON)
		result := dbFilter.CheckFilter(tc.dbName, tc.collName)
		assert.Equal(t, tc.expBool, result)
	}
}

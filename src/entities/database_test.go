package entities

import (
	"testing"

	"github.com/globalsign/mgo/bson"

	"github.com/stretchr/testify/mock"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/filter"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
)

func Test_databaseCollector_GetEntity(t *testing.T) {
	dc := getTestDatabaseCollector()

	e, err := dc.GetEntity()
	assert.NoError(t, err)
	assert.Equal(t, "testDatabase", e.Metadata.Name)
	assert.Equal(t, "database", e.Metadata.Namespace)
}

func Test_databaseCollector_GetEntity_Error(t *testing.T) {
	dc := getBadTestDatabaseCollector()
	e, err := dc.GetEntity()
	assert.Error(t, err)
	assert.Nil(t, e)
}

func Test_databaseCollector_CollectInventory(t *testing.T) {
	cc := getTestDatabaseCollector()
	assert.NotPanics(t, func() {
		cc.CollectInventory()
	})
}

func Test_databaseCollector_CollectMetrics(t *testing.T) {
	i, _ := integration.New("test", "0.0.1")
	dc := databaseCollector{
		defaultCollector{
			"testDatabase",
			i,
			test.FakeSession{},
		},
	}

	dc.CollectMetrics()
}

func Test_databaseCollector_CollectMetrics_Error(t *testing.T) {
	dc := getBadTestDatabaseCollector()
	assert.NotPanics(t, func() {
		dc.CollectMetrics()
	})
}

func TestGetDatabases(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	testFilter, _ := filter.ParseFilters("")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", Cmd{"listDatabases": 1}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"databases": [
					{
						"name": "db1"
					}, {
						"name": "db2"
					}
				]
			}`), result)
			assert.NoError(t, err)
		}).
		Once()
	expectedDBs := []string{"db1", "db2"}

	databases, err := GetDatabases(mockSession, testIntegration, testFilter)
	mockSession.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotEmpty(t, databases)
	assert.Equal(t, len(expectedDBs), len(databases))
	for i, collector := range databases {
		session, err := collector.GetSession()
		assert.NoError(t, err)
		assert.Equal(t, mockSession, session)
		assert.Equal(t, testIntegration, collector.GetIntegration())
		assert.Equal(t, expectedDBs[i], collector.GetName())
	}
}

func TestGetDatabases_Error(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	testFilter, _ := filter.ParseFilters("")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", Cmd{"listDatabases": 1}, mock.Anything).
		Return(assert.AnError).
		Once()

	databases, err := GetDatabases(mockSession, testIntegration, testFilter)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Nil(t, databases)
}

func getTestDatabaseCollector() *databaseCollector {
	i, _ := integration.New("testIntegration", "testVersion")
	return &databaseCollector{
		defaultCollector{
			"testDatabase",
			i,
			test.FakeSession{},
		},
	}
}

func getBadTestDatabaseCollector() *databaseCollector {
	return &databaseCollector{
		defaultCollector{
			"testDatabase",
			nil,
			nil,
		},
	}
}

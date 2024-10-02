package entities

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/mock"

	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
)

func Test_mongosCollector_GetEntity(t *testing.T) {
	mc := getTestMongosCollector()

	e, err := mc.GetEntity()
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, "testMongos", e.Metadata.Name)
	assert.Equal(t, "mo-mongos", e.Metadata.Namespace)
}

func Test_mongosCollector_GetEntity_Error(t *testing.T) {
	mc := getBadTestMongosCollector()

	e, err := mc.GetEntity()
	assert.Error(t, err)
	assert.Nil(t, e)
}

func Test_mongosCollector_CollectInventory(t *testing.T) {
	mc := getTestMongosCollector()
	mc.CollectInventory()
	e, err := mc.GetEntity()
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, test.ExpectedInventory, e.Inventory.Items())
}

func Test_mongosCollector_CollectInventory_Error(t *testing.T) {
	mc := getBadTestMongosCollector()
	mc.CollectInventory()
	assert.NotPanics(t, func() {
		mc.CollectInventory()
	})
}

func Test_mongosCollector_CollectMetrics(t *testing.T) {
	mc := getTestMongosCollector()
	mc.CollectMetrics()
	e, err := mc.GetEntity()
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.NotEmpty(t, e.Metrics)
}

func Test_mongosCollector_CollectMetrics_Error(t *testing.T) {
	mc := getBadTestMongosCollector()
	assert.NotPanics(t, func() {
		mc.CollectMetrics()
	})
}

func TestGetMongoses(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.On("New", "host1", "27017").Return(mockSession, nil).Once()
	mockSession.On("New", "host2", "1234").Return(mockSession, nil).Once()
	mockSession.On("New", "host3", "666").Return(nil, assert.AnError).Once()
	mockSession.MockDatabase("config", 1).
		MockCollection("mongos", 1).
		On("FindAll", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(0)
			err := bson.UnmarshalJSON([]byte(`[
				{ "_id": "host1:27017" },
				{ "_id": "host2:1234" },
				{ "_id": "host3:666" }
			]`), result)
			assert.NoError(t, err)
		}).
		Once()

	collectors, err := GetMongoses(mockSession, testIntegration)
	expectedHosts := []string{"host1:27017", "host2:1234"}

	mockSession.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, collectors, len(expectedHosts))
	for i, coll := range collectors {
		session, err := coll.GetSession()
		assert.NoError(t, err)
		assert.Equal(t, mockSession, session)
		assert.Equal(t, testIntegration, coll.GetIntegration())
		assert.Equal(t, expectedHosts[i], coll.GetName())
	}
}

func TestGetMongoses_Error(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("config", 1).MockCollection("mongos", 1).
		On("FindAll", mock.Anything).
		Return(assert.AnError).
		Once()

	collectors, err := GetMongoses(mockSession, testIntegration)

	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Nil(t, collectors)
}

func getTestMongosCollector() *mongosCollector {
	i, _ := integration.New("testIntegration", "testVersion")
	return &mongosCollector{
		hostCollector{
			defaultCollector{
				"testMongos",
				i,
				test.FakeSession{},
				nil,
			},
		},
	}
}

func getBadTestMongosCollector() *mongosCollector {
	return &mongosCollector{
		hostCollector{
			defaultCollector{
				"testMongos",
				nil,
				nil,
				nil,
			},
		},
	}
}

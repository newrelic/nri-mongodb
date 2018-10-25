package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
)

func Test_mongodCollector_GetEntity(t *testing.T) {
	mc := getTestMongodCollector()

	e, err := mc.GetEntity()
	assert.NoError(t, err)
	assert.Equal(t, "testMongod", e.Metadata.Name)
	assert.Equal(t, "mongod", e.Metadata.Namespace)
}

func Test_mongodCollector_GetEntity_Error(t *testing.T) {
	mc := getBadTestMongodCollector()

	e, err := mc.GetEntity()
	assert.Error(t, err)
	assert.Nil(t, e)
}

func Test_mongodCollector_CollectInventory(t *testing.T) {
	mc := getTestMongodCollector()
	mc.CollectInventory()
	e, err := mc.GetEntity()
	assert.NoError(t, err)
	assert.Equal(t, test.ExpectedInventory, e.Inventory.Items())
}

func Test_mongodCollector_CollectInventory_Error(t *testing.T) {
	mc := getBadTestMongodCollector()
	assert.NotPanics(t, func() {
		mc.CollectInventory()
	})
}

func Test_mongodCollector_CollectMetrics(t *testing.T) {
	mc := getTestMongodCollector()
	mc.CollectMetrics()
	e, err := mc.GetEntity()
	assert.NoError(t, err)
	assert.NotEmpty(t, e.Metrics)
}

func Test_mongodCollector_CollectMetrics_Error(t *testing.T) {
	mc := getBadTestMongodCollector()
	assert.NotPanics(t, func() {
		mc.CollectMetrics()
	})
}

func TestGetMongods(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.On("New", "host1", "1234").Return(mockSession, nil).Once()
	mockSession.On("New", "host2", "4321").Return(mockSession, nil).Once()
	mockSession.On("New", "host3", "666").Return(nil, assert.AnError).Once()

	expectedHosts := []string{"host1:1234", "host2:4321"}
	collectors, err := GetMongods(mockSession, "rs1/host1:1234,host2:4321,host3:666", testIntegration)

	mockSession.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedHosts), len(collectors))

	for i, collector := range collectors {
		session, err := collector.GetSession()
		assert.NoError(t, err)
		assert.Equal(t, testIntegration, collector.GetIntegration())
		assert.Equal(t, mockSession, session)
		assert.Equal(t, expectedHosts[i], collector.GetName())
	}
}

func getTestMongodCollector() *mongodCollector {
	i, _ := integration.New("testIntegration", "testVersion")
	return &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				i,
				test.FakeSession{},
			},
		},
	}
}

func getBadTestMongodCollector() *mongodCollector {
	return &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				nil,
				nil,
			},
		},
	}
}

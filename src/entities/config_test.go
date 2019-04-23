package entities

import (
	"testing"

	"github.com/globalsign/mgo/bson"

	"github.com/stretchr/testify/mock"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
)

func Test_configCollector_GetEntity(t *testing.T) {
	cc := getTestConfigCollector()
	e, err := cc.GetEntity()
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, "testConfig", e.Metadata.Name)
	assert.Equal(t, "mo-config", e.Metadata.Namespace)
}

func Test_configCollector_GetEntity_Error(t *testing.T) {
	cc := getBadTestConfigCollector()
	e, err := cc.GetEntity()
	assert.Error(t, err)
	assert.Nil(t, e)
}

func Test_configCollector_CollectInventory(t *testing.T) {
	cc := getTestConfigCollector()
	assert.NotPanics(t, func() {
		cc.CollectInventory()
	})
}

func Test_configCollector_CollectInventory_Error(t *testing.T) {
	cc := getBadTestConfigCollector()
	assert.NotPanics(t, func() {
		cc.CollectInventory()
	})
}

func Test_configCollector_CollectMetrics(t *testing.T) {
	cc := getTestConfigCollector()
	assert.NotPanics(t, func() {
		cc.CollectMetrics()
	})
}

func Test_configCollector_CollectMetrics_Error(t *testing.T) {
	cc := getBadTestConfigCollector()
	assert.NotPanics(t, func() {
		cc.CollectMetrics()
	})
}

func TestGetConfigServers(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.On("New", "config1", "27017").Return(mockSession, nil).Once()
	mockSession.On("New", "config2", "27017").Return(mockSession, nil).Once()
	mockSession.On("New", "config3", "666").Return(nil, assert.AnError).Once()
	mockSession.MockDatabase("admin", 1).
		On("Run", "getShardMap", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			bson.UnmarshalJSON([]byte(`{
				"map": {
					"config": "rs1config/config1:27017,config2:27017,config3:666"
				}
			}`), result)
		}).
		Once()
	expectedConfigServers := []struct {
		name string
		port string
	}{
		{"config1:27017", "27017"},
		{"config2:27017", "27017"},
	}

	collectors, err := GetConfigServers(mockSession, testIntegration)
	mockSession.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotEmpty(t, collectors)
	assert.Equal(t, len(expectedConfigServers), len(collectors))
	for i, collector := range collectors {
		session, err := collector.GetSession()
		assert.NoError(t, err)
		assert.Equal(t, mockSession, session)
		assert.Equal(t, testIntegration, collector.GetIntegration())
		assert.Equal(t, expectedConfigServers[i].name, collector.GetName())
	}
}

func TestGetConfigServers_ErrorRunningCommand(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", "getShardMap", mock.Anything).
		Return(assert.AnError).
		Once()

	collectors, err := GetConfigServers(mockSession, testIntegration)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, collectors)
}

func TestGetConfigServers_EmptyConfigHost(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", "getShardMap", mock.Anything).
		Return(nil).
		Once()

	collectors, err := GetConfigServers(mockSession, testIntegration)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, collectors)
}

func TestGetConfigServers_CannotConnectToSession(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.On("New", "config1", "666").Return(nil, assert.AnError).Once()
	mockSession.MockDatabase("admin", 1).
		On("Run", "getShardMap", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			bson.UnmarshalJSON([]byte(`{
				"map": {
					"config": "rs1config/config1:666"
				}
			}`), result)
		}).
		Once()

	collectors, err := GetConfigServers(mockSession, testIntegration)
	mockSession.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Empty(t, collectors)
}

func getTestConfigCollector() *configCollector {
	i, _ := integration.New("testIntegration", "testVersion")

	return &configCollector{
		hostCollector{
			defaultCollector{
				"testConfig",
				i,
				test.FakeSession{},
			},
		},
	}
}

func getBadTestConfigCollector() *configCollector {
	return &configCollector{
		hostCollector{
			defaultCollector{
				"testConfig",
				nil,
				nil,
			},
		},
	}
}

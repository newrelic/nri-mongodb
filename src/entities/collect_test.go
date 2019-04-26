package entities

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_collectServerStatus(t *testing.T) {
	c := getTestMongodCollector()

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset", metric.Attribute{Key: "key", Value: "value"})

	err := collectServerStatus(c, ms)
	if err != nil {
		t.Error(err)
	}
	expected := map[string]interface{}{
		"asserts.regularPerSecond":   float64(0),
		"asserts.warningPerSecond":   float64(0),
		"asserts.messagesPerSecond":  float64(0),
		"asserts.userPerSecond":      float64(0),
		"asserts.rolloversPerSecond": float64(0),
		"key":                "value",
		"event_type":         "testmetricset",
		"reportingEntityKey": "mo-mongod:testhost:1234:clustername=",
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectServerStatus_MissingSession(t *testing.T) {
	c := getTestMongodCollector()
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	c.session = nil
	expectedCount := len(ms.Metrics)

	err := collectServerStatus(c, ms)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount) // 1 for the eventType
}

func Test_collectServerStatus_CommandError(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", Cmd{"serverStatus": 1}, mock.Anything).
		Return(assert.AnError).
		Once()
	mockSession.MockDatabase("admin", 1).
		On("Run", Cmd{"isMaster": 1}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"isMaster": true,
				"ok": 1
			}`), result)
			assert.NoError(t, err)
		}).
		Once()

	c := getTestMongodCollector()
	c.session = mockSession
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectServerStatus(c, ms)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectIsMaster(t *testing.T) {
	c := getTestMongodCollector()

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	_, err := collectIsMaster(c, ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"replset.isMaster":    float64(1),
		"replset.isSecondary": float64(1),
		"event_type":          "testmetricset",
		"reportingEntityKey":  "mo-mongod:testhost:1234:clustername=",
	}

	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func TestIsStandaloneInstance_True(t *testing.T) {
	mockSession := new(test.MockSession)
	adminDB := mockSession.MockDatabase("admin", 1)
	adminDB.On("Run", Cmd{"isMaster": 1}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"isMaster": true,
				"ok": 1
			}`), result)
			assert.NoError(t, err)
		}).
		Once()

	result, err := IsStandaloneInstance(mockSession)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, result)
}

func TestIsStandaloneInstance_False(t *testing.T) {
	mockSession := new(test.MockSession)
	adminDB := mockSession.MockDatabase("admin", 1)
	adminDB.On("Run", Cmd{"isMaster": 1}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"isMaster": true,
				"msg": "isdbgrid",
				"ok": 1
			}`), result)
			assert.NoError(t, err)
		}).
		Once()

	result, err := IsStandaloneInstance(mockSession)
	if err != nil {
		t.Error(err)
	}

	assert.False(t, result)
}

func Test_collectIsMaster_MissingSession(t *testing.T) {
	c := getTestMongodCollector()
	c.session = nil
	e, err := c.GetEntity()
	assert.Error(t, err)
	assert.Nil(t, e)
}

func Test_collectIsMaster_CommandError(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", Cmd{"isMaster": 1}, mock.Anything).
		Return(assert.AnError).
		Once()

	c := getTestMongodCollector()
	c.session = mockSession
	e, err := c.GetEntity()
	assert.Error(t, err)
	assert.Nil(t, e)
}

func Test_collectReplGetStatus(t *testing.T) {
	c := getTestMongodCollector()

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := collectReplGetStatus(c, "mdb-rh7-rs1-a1.bluemedora.localnet:27017", ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"replset.health":               float64(1),
		"replset.state":                "SECONDARY",
		"replset.uptimeInMilliseconds": float64(758657),
		"event_type":                   "testmetricset",
		"reportingEntityKey":           "mo-mongod:testhost:1234:clustername=",
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectReplGetConfig(t *testing.T) {
	c := getTestMongodCollector()

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := collectReplGetConfig(c, "mdb-rh7-rs1-a1.bluemedora.localnet:27017", ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"replset.isArbiter":    float64(0),
		"replset.isHidden":     float64(0),
		"replset.priority":     float64(10),
		"replset.votes":        float64(20),
		"replset.voteFraction": float64(1),
		"event_type":           "testmetricset",
		"reportingEntityKey":   "mo-mongod:testhost:1234:clustername=",
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectTop(t *testing.T) {
	c := getTestMongodCollector()

	e, _ := c.GetEntity()

	err := collectTop(c)
	if err != nil {
		t.Error(err)
	}
	expected := map[string]interface{}{
		"usage.totalInMillisecondsPerSecond":     float64(0),
		"usage.totalPerSecond":                   float64(0),
		"usage.writeLockPerSecond":               float64(0),
		"event_type":                             "MongodTopSample",
		"displayName":                            "testMongod",
		"database":                               "records",
		"collection":                             "users",
		"entityName":                             "mo-mongod:testMongod",
		"usage.readLockInMillisecondsPerSecond":  float64(0),
		"usage.readLockPerSecond":                float64(0),
		"usage.writeLockInMillisecondsPerSecond": float64(0),
		"reportingEntityKey":                     "mo-mongod:testhost:1234:clustername=",
	}
	actual := e.Metrics[0].Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectCollStats(t *testing.T) {
	c := getTestCollectionCollector()

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")

	err := collectCollStats(c, ms)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"collection.sizeInBytes":       float64(2157),
		"collection.avgObjSizeInBytes": float64(719),
		"collection.count":             float64(3),
		"collection.capped":            float64(0),
		"event_type":                   "test",
		"reportingEntityKey":           "mo-mongod:testhost:1234:clustername=",
	}
	assert.Equal(t, expected, ms.Metrics)
}

func Test_collectCollStats_SkipSystemCollection(t *testing.T) {
	c := getTestCollectionCollector()
	c.name = "system.admin"
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectCollStats(c, ms)
	assert.NoError(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectDbStats(t *testing.T) {
	c := getTestDatabaseCollector()

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")

	err := collectDbStats(c, ms)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"stats.objects":        float64(5),
		"stats.storageInBytes": float64(7),
		"stats.indexInBytes":   float64(8),
		"stats.indexes":        float64(4),
		"stats.dataInBytes":    float64(6),
		"event_type":           "test",
		"reportingEntityKey":   "mo-mongod:testhost:1234:clustername=",
	}
	assert.Equal(t, expected, ms.Metrics)
}

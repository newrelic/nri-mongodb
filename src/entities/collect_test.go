package entities

import (
	"testing"

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
		"key":        "value",
		"event_type": "testmetricset",
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
		On("Run", cmd{"serverStatus": 1}, mock.Anything).
		Return(assert.AnError).
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
	}

	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectIsMaster_MissingSession(t *testing.T) {
	c := getTestMongodCollector()
	c.session = nil
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	isMaster, err := collectIsMaster(c, ms)
	assert.Error(t, err)
	assert.False(t, isMaster)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectIsMaster_CommandError(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", cmd{"isMaster": 1}, mock.Anything).
		Return(assert.AnError).
		Once()

	c := getTestMongodCollector()
	c.session = mockSession
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	isMaster, err := collectIsMaster(c, ms)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.False(t, isMaster)
	assert.Len(t, ms.Metrics, expectedCount)
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
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectReplGetStatus_MissingSession(t *testing.T) {
	c := getTestMongodCollector()
	c.session = nil
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectReplGetStatus(c, "", ms)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectReplGetStatus_CommandError(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", cmd{"replSetGetStatus": 1}, mock.Anything).
		Return(assert.AnError).
		Once()

	c := getTestMongodCollector()
	c.session = mockSession
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectReplGetStatus(c, "", ms)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
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
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectReplGetConfig_MissingSession(t *testing.T) {
	c := getTestMongodCollector()
	c.session = nil
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectReplGetConfig(c, "", ms)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectReplGetConfig_CommandError(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", cmd{"replSetGetConfig": 1}, mock.Anything).
		Return(assert.AnError).
		Once()

	c := getTestMongodCollector()
	c.session = mockSession
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectReplGetConfig(c, "", ms)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectTop(t *testing.T) {
	c := getTestMongodCollector()

	e, _ := c.GetEntity()

	err := collectTop(c)
	if err != nil {
		t.Error(err)
	}
	expected := map[string]interface{}{
		"usage.totalInMilliseconds":     float64(305277),
		"usage.totalPerSecond":          float64(0),
		"usage.writeLockPerSecond":      float64(0),
		"event_type":                    "MongodTopSample",
		"displayName":                   "testMongod",
		"database":                      "records",
		"collection":                    "users",
		"entityName":                    "mongod:testMongod",
		"usage.readLockInMilliseconds":  float64(305123),
		"usage.readLockPerSecond":       float64(0),
		"usage.writeLockInMilliseconds": float64(13),
	}
	actual := e.Metrics[0].Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectTop_MissingSession(t *testing.T) {
	c := getTestMongodCollector()
	c.session = nil
	e, _ := c.GetEntity()

	err := collectTop(c)
	assert.Error(t, err)
	assert.Empty(t, e.Metrics)
}

func Test_collectTop_MissingIntegration(t *testing.T) {
	c := getTestMongodCollector()
	c.integration = nil

	err := collectTop(c)
	assert.Error(t, err)
}

func Test_collectTop_CommandError(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("admin", 1).
		On("Run", cmd{"top": 1}, mock.Anything).
		Return(assert.AnError).
		Once()

	c := getTestMongodCollector()
	c.session = mockSession
	e, _ := c.GetEntity()

	err := collectTop(c)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Empty(t, e.Metrics)
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

func Test_collectCollStats_MissingSession(t *testing.T) {
	c := getTestCollectionCollector()
	c.session = nil
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectCollStats(c, ms)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

func Test_collectCollStats_CommandError(t *testing.T) {
	c := getTestCollectionCollector()
	mockSession := new(test.MockSession)
	mockSession.MockDatabase(c.db, 1).
		On("Run", cmd{"collStats": c.name}, mock.Anything).
		Return(assert.AnError).
		Once()

	c.session = mockSession // mocks.Session
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectCollStats(c, ms)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
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
	}
	assert.Equal(t, expected, ms.Metrics)
}

func Test_collectDbStats_CommandError(t *testing.T) {
	c := getTestDatabaseCollector()
	mockSession := new(test.MockSession)
	mockSession.MockDatabase(c.name, 1).
		On("Run", cmd{"dbStats": 1}, mock.Anything).
		Return(assert.AnError).
		Once()
	c.session = mockSession
	e, _ := c.GetEntity()
	ms := e.NewMetricSet("test")
	expectedCount := len(ms.Metrics)

	err := collectDbStats(c, ms)
	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, ms.Metrics, expectedCount)
}

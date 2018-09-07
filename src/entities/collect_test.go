package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
)

func TestCollectServerStatus(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				i,
				test.MockSession{},
			},
		},
	}

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

func Test_collectIsMaster(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				i,
				test.MockSession{},
			},
		},
	}

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

func Test_collectReplGetStatus(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				i,
				test.MockSession{},
			},
		},
	}

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

func Test_collectReplGetConfig(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				i,
				test.MockSession{},
			},
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := collectReplGetConfig(c, "mdb-rh7-rs1-a1.bluemedora.localnet:27017", ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"replset.isArbiter": float64(0),
		"replset.isHidden":  float64(0),
		"replset.priority":  float64(10),
		"replset.votes":     float64(20),
		"event_type":        "testmetricset",
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

func Test_collectTop(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &mongodCollector{
		hostCollector{
			defaultCollector{
				"testMongod",
				i,
				test.MockSession{},
			},
		},
	}

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

func Test_collectCollStats(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &collectionCollector{
		defaultCollector{
			"testMongod",
			i,
			test.MockSession{},
		},
		"testDB",
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := collectCollStats(c, ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"collection.sizeInBytes":       float64(2157),
		"collection.avgObjSizeInBytes": float64(719),
		"collection.count":             float64(3),
		"collection.capped":            float64(0),
		"event_type":                   "testmetricset",
	}
	actual := ms.Metrics

	assert.Equal(t, expected, actual)
}

func Test_collectDbStats(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := &databaseCollector{
		defaultCollector{
			"testMongod",
			i,
			test.MockSession{},
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := collectDbStats(c, ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"stats.objects":        float64(5),
		"stats.storageInBytes": float64(7),
		"stats.indexInBytes":   float64(8),
		"stats.indexes":        float64(4),
		"stats.dataInBytes":    float64(6),
		"event_type":           "testmetricset",
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)
}

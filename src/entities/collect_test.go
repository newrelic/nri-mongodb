package entities

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tonnerre/golang-pretty"

	"github.com/newrelic/infra-integrations-sdk/data/metric"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func TestCollectServerStatus(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := MongodCollector{
		HostCollector{
			DefaultCollector{
				Session:     test.MockSession{},
				Integration: i,
			},
			"testMongod",
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset", metric.Attribute{Key: "key", Value: "value"})

	err := CollectServerStatus(c, ms)
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
	fmt.Println(pretty.Diff(actual, expected))
	assert.True(t, reflect.DeepEqual(actual, expected))
}

func TestCollectIsMaster(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := MongodCollector{
		HostCollector{
			DefaultCollector{
				Session:     test.MockSession{},
				Integration: i,
			},
			"testMongod",
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	_, err := CollectIsMaster(c, ms)
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
	// 	assert.True(t, reflect.DeepEqual(actual, expected))
}

func TestCollectReplGetStatus(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := MongodCollector{
		HostCollector{
			DefaultCollector{
				Session:     test.MockSession{},
				Integration: i,
			},
			"testMongod",
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := CollectReplGetStatus(c, "mdb-rh7-rs1-a1.bluemedora.localnet:27017", ms)
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

func TestCollectReplGetConfig(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := MongodCollector{
		HostCollector{
			DefaultCollector{
				Session:     test.MockSession{},
				Integration: i,
			},
			"testMongod",
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := CollectReplGetConfig(c, "mdb-rh7-rs1-a1.bluemedora.localnet:27017", ms)
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

func TestCollectTop(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := MongodCollector{
		HostCollector{
			DefaultCollector{
				Session:     test.MockSession{},
				Integration: i,
			},
			"testMongod",
		},
	}

	e, _ := c.GetEntity()

	err := CollectTop(c)
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

	// for i, ms := range e.Metrics {
	// 	assert.Equal(t, "", ms)
	// 	print(i)
	// }

}

func TestCollectCollStats(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := CollectionCollector{
		DefaultCollector{
			Session:     test.MockSession{},
			Integration: i,
		},
		"testMongod",
		"testDB",
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := CollectCollStats(c, ms)
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

func TestCollectDbStats(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := DatabaseCollector{
		DefaultCollector{
			Session:     test.MockSession{},
			Integration: i,
		},
		"testMongod",
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := CollectDbStats(c, ms)
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

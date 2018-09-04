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
		"asserts.regularPerSecond":   0.0,
		"asserts.warningPerSecond":   0.0,
		"asserts.messagesPerSecond":  0.0,
		"asserts.userPerSecond":      0.0,
		"asserts.rolloversPerSecond": 0.0,
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
		"replset.isMaster":    true,
		"replset.isSecondary": true,
		"key":        "value",
		"event_type": "testmetricset",
	}

	actual := ms.Metrics
	assert.Equal(t, expected, actual)
	// 	assert.True(t, reflect.DeepEqual(actual, expected))
}

func TestCollectReplSetMetrics(t *testing.T) {

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

	err := CollectReplSetMetrics(c, ms)
	if err != nil {
		t.Error(err)
	}

	expected := map[string]interface{}{
		"members": []map[string]interface{}{
			{
				"name":                         "mdb-rh7-rs1-a1.bluemedora.localnet:27017",
				"replset.health":               0.0,
				"replset.stateStr":             "SECONDARY",
				"replset.uptimeInMilliseconds": 0.0,
			},
		},
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

	err := CollectTop(c)
	if err != nil {
		t.Error(err)
	}

	//not entirely sure what actual should be, since we never create ms like in the other tests

	// expected := map[string]interface{}{
	// 	"totals": map[string]interface{}{
	// 		"total": map[string]interface{}{
	// 			"usage.totalInMilliseconds": 0,
	// 			"usage.totalPerSecond":      0,
	// 		},
	// 	},
	// }
	// actual := ms.Metrics
	// assert.Equal(t, expected, actual)
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
		"collection.sizeInBytes":       2157.0,
		"collection.avgObjSizeInBytes": 719.0,
		"collection.count":             3.0,
		"collection.capped":            0.0,
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
		"stats.objects":        5.0,
		"stats.storageInBytes": 7.0,
		"stats.indexInBytes":   8.0,
		"stats.indexes":        4.0,
		"stats.dataInBytes":    6.0,
		"event_type":           "testmetricset",
	}
	actual := ms.Metrics
	assert.Equal(t, expected, actual)

}

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/entities"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartCollectorWorkerPool(t *testing.T) {
	numWorkers := 10
	var wg sync.WaitGroup
	entitiesChan := StartCollectorWorkerPool(numWorkers, &wg)
	close(entitiesChan)

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return
	case <-time.After(time.Second):
		assert.FailNow(t, "Wait group close timed out")
	}
}

type testCollector struct {
	name        string
	integration *integration.Integration
	session     connection.Session
}

func (t *testCollector) GetEntity() (*integration.Entity, error) {
	if t.integration != nil {
		return t.integration.Entity(t.name, "test")
	}

	return nil, assert.AnError
}

func (t *testCollector) GetName() string {
	return t.name
}

func (t *testCollector) GetIntegration() *integration.Integration {
	return t.integration
}

func (t *testCollector) GetSession() (connection.Session, error) {
	return t.session, nil
}

func (t *testCollector) CollectInventory() {
	e, _ := t.GetEntity()
	e.SetInventoryItem("testCategory", "testItem", "testValue")
	return
}

func (t *testCollector) CollectMetrics() {
	e, _ := t.GetEntity()

	ms := e.NewMetricSet("testSample")
	ms.SetMetric("testMetric", 1, metric.GAUGE)
	return
}

func Test_collectorWorker(t *testing.T) {
	collectorChan := make(chan entities.Collector)
	var wg sync.WaitGroup
	i, _ := integration.New("testIntegration", "testVersion")

	wg.Add(1)
	go collectorWorker(collectorChan, &wg)

	collectorChan <- &testCollector{
		"testName",
		i,
		test.FakeSession{},
	}
	close(collectorChan)

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		assert.Len(t, i.Entities, 1, "Expected one entity")
		assert.Len(t, i.Entities[0].Metrics[0].Metrics, 2, "Expected one metric in the set")
		assert.Len(t, i.Entities[0].Inventory.Items(), 1, "Expected one inventory item")
	case <-time.After(time.Second):
		assert.FailNow(t, "Collector worker took too long to close.")
	}
}

func TestFeedWorkerPool(t *testing.T) {
	mockSession := new(test.MockSession)
	mockSession.On("New", "config1", "27017").Return(mockSession, nil).Once()
	mockSession.On("New", "mongos1", "27017").Return(mockSession, nil).Once()
	mockSession.On("New", "shard1", "27017").Return(mockSession, nil).Once()

	configDB := mockSession.MockDatabase("config", 3)
	configDB.MockCollection("mongos", 2).
		On("FindAll", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(0)
			err := bson.UnmarshalJSON([]byte(`[
				{ "_id": "mongos1:27017" },
			]`), result)
			assert.NoError(t, err)
		})
	configDB.MockCollection("shards", 1).
		On("FindAll", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(0)
			err := bson.UnmarshalJSON([]byte(`[
				{ "_id": "rs1", "host": "shard1" },
			]`), result)
			assert.NoError(t, err)
		})

	adminDB := mockSession.MockDatabase("admin", 3)

	adminDB.On("Run", entities.Cmd{"isMaster": 1}, mock.Anything).
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

	adminDB.On("Run", entities.Cmd{"listDatabases": 1}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"databases": [
					{
						"name": "database1"
					}
				]
			}`), result)
			assert.NoError(t, err)
		}).
		Once()

	adminDB.On("Run", "getShardMap", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"map": {
					"config": "rs1config/config1:27017"
				}
			}`), result)
			assert.NoError(t, err)
		}).
		Once()

	mockSession.MockDatabase("database1", 1).
		On("CollectionNames").
		Return([]string{"collection1"}, nil).
		Once()

	collChan := make(chan entities.Collector)
	i, _ := integration.New("test", "0.0.0")

	go FeedWorkerPool(mockSession, collChan, i)

	wgDone := make(chan struct{})
	var collectors []entities.Collector
	go func() {
		for {
			coll, ok := <-collChan
			if !ok {
				break
			} else {
				collectors = append(collectors, coll)
			}
		}
		close(wgDone)
	}()

	expectedCollectorNames := map[string]bool{
		"database1":     true,
		"config1:27017": true,
		"mongos1:27017": true,
		"mongos1":       true,
		"shard1:27017":  true,
		"collection1":   true,
	}

	select {
	case <-wgDone:
		mockSession.AssertExpectations(t)
		assert.Len(t, collectors, len(expectedCollectorNames))
		for _, coll := range collectors {
			_, ok := expectedCollectorNames[coll.GetName()]
			assert.True(t, ok, "Expected collector name is missing: %s", coll.GetName())
			session, err := coll.GetSession()
			assert.NoError(t, err)
			assert.Equal(t, mockSession, session)
			assert.Equal(t, i, coll.GetIntegration())
			e, err := coll.GetEntity()
			assert.NoError(t, err)
			assert.NotNil(t, e)
		}
	case <-time.After(time.Second):
		assert.FailNow(t, "Timed out waiting for Mongoses")
	}
}

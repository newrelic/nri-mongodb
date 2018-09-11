package main

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/entities"
	"github.com/newrelic/nri-mongodb/src/test"
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
		t.Error("wait group close timed out")
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

	return nil, errors.New("nil integration")
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
		test.MockSession{},
	}
	close(collectorChan)

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		if len(i.Entities) != 1 {
			t.Errorf("expected one entity, got %d", len(i.Entities))
		}
		if len(i.Entities[0].Metrics[0].Metrics) != 2 {
			t.Errorf("expected one metric in the set, got %d", len(i.Entities[0].Metrics[0].Metrics))
		}

		if len(i.Entities[0].Inventory.Items()) != 1 {
			t.Errorf("expected one inventory item, got %d", len(i.Entities[0].Inventory.Items()))
		}
	case <-time.After(time.Second):
		t.Error("collector worker took too long to close")
	}
}

package main

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/entities"
	"sync"
	"testing"
	"time"
)

func TestStartCollectorWorkerPool(t *testing.T) {
	numWorkers := 10
	var wg sync.WaitGroup
	i, _ := integration.New("testIntegration", "testVersion")

	entitiesChan := StartCollectorWorkerPool(numWorkers, &wg, i)
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
	name string
}

func (t testCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(t.name, "testEntity")
}

func (t testCollector) CollectInventory(e *integration.Entity) {
	e.SetInventoryItem("testCategory", "testItem", "testValue")
	return
}

func (t testCollector) CollectMetrics(e *integration.Entity) {
	ms := e.NewMetricSet("testSample")
	ms.SetMetric("testMetric", 1, metric.GAUGE)
	return
}

func Test_collectorWorker(t *testing.T) {
	collectorChan := make(chan entities.Collector)
	var wg sync.WaitGroup
	i, _ := integration.New("testIntegration", "testVersion")

	wg.Add(1)
	go collectorWorker(collectorChan, &wg, i)

	collectorChan <- testCollector{name: "testName"}
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

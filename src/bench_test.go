package main

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/newrelic/nri-mongodb/src/metrics"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/stretchr/testify/mock"
	"strconv"
	"sync"
	"testing"
)

func Benchmark1(b *testing.B) {
	benchmarkDatabases(1, b)
}

func Benchmark10(b *testing.B) {
	benchmarkDatabases(10, b)
}

func Benchmark100(b *testing.B) {
	benchmarkDatabases(100, b)
}

func Benchmark1000(b *testing.B) {
	benchmarkDatabases(1000, b)
}

func Benchmark10000(b *testing.B) {
	benchmarkDatabases(10000, b)
}

func Benchmark100000(b *testing.B) {
	benchmarkDatabases(100000, b)
}

func benchmarkDatabases(size int, b *testing.B) {
	for n := 0; n < b.N; n++ {
    mockSession := &test.MockSession{}
		var wg sync.WaitGroup
		testIntegration, _ := integration.New("Test", "0.0.1")
		collectorChan := StartCollectorWorkerPool(10, &wg)
    for i := 0; i < size; i++ {
      mdb := &MockDatabaseCollector{
        name: "db" + strconv.Itoa(i),
        integration: testIntegration,
        session: mockSession,
      }

      collectorChan <- mdb
    }
	}
}

type MockDatabaseCollector struct {
	mock.Mock
	name        string
	integration *integration.Integration
	session     connection.Session
}


func (m *MockDatabaseCollector) CollectMetrics() {
  e, _ := m.GetEntity()
  i := 1
  f := false
  stats := metrics.DbStats {
    Objects: &i,
    StorageSize: &i,
    IndexSize: &i,
    Indexes: &i,
    DataSize: &i,
  }

  coll := metrics.CollStats {
    Size: &i,
    AvgObjSize: &i,
    Count: &i,
    Capped: &f,
    Max: &i,
    MaxSize: &i,
    StorageSize:&i, 
    Nindexes: &i,
    IndexSizes: &map[string]int{"test":1},
  }

  dms := e.NewMetricSet("MongoDatabaseSample")
  dms.MarshalMetrics(stats)

  cms := e.NewMetricSet("MongoCollectionSample")
  cms.MarshalMetrics(coll)
}

  
func (m *MockDatabaseCollector) CollectInventory() {
  return
}

func (m *MockDatabaseCollector) GetName() string {
	return m.name
}

func (m *MockDatabaseCollector) GetEntity() (*integration.Entity, error) {
	return m.GetIntegration().Entity(m.name, "database")
}

func (m *MockDatabaseCollector) GetIntegration() (*integration.Integration){
	return m.integration
}

func (m *MockDatabaseCollector) GetSession() (connection.Session, error){
  return m.session, nil
}

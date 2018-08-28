package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_CollectionCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")

	cc := CollectionCollector{
		DefaultCollector{
			Integration: i,
			Session:     test.MockSession{},
		},
		"testCollector",
		"testDB",
	}

	e, err := cc.GetEntity()
	if err != nil {
		t.Error(err)
	}

	if e.Metadata.Name != "testCollector" {
		t.Errorf("Expected entity name testCollector, got %s", e.Metadata.Name)
	}

	if e.Metadata.Namespace != "collection" {
		t.Errorf("Expected entity namespace collection, got %s", e.Metadata.Namespace)
	}

}

/*
 *func Test_CollectionCollector_CollectMetrics(t *testing.T) {
 *  cc := CollectionCollector{
 *    Name: "testCollector",
 *    DB:   "testDB",
 *  }
 *
 *  i, _ := integration.New("testIntegration", "testVersion")
 *
 *  e, err := cc.GetEntity(i)
 *  if err != nil {
 *    t.Error(err)
 *  }
 *
 *}
 */

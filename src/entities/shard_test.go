package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_ShardCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")
	cc := ShardCollector{
		DefaultCollector{
			Integration: i,
			Session:     test.MockSession{},
		},
		"testCollector",
		"testHost",
	}

	e, err := cc.GetEntity()
	if err != nil {
		t.Error(err)
	}

	if e.Metadata.Name != "testHost" {
		t.Errorf("Expected entity name testCollector, got %s", e.Metadata.Name)
	}

	if e.Metadata.Namespace != "shard" {
		t.Errorf("Expected entity namespace shard, got %s", e.Metadata.Namespace)
	}

}

func Test_ShardCollector_CollectMetrics(t *testing.T) {
	i, _ := integration.New("test", "0.0.1")
	cc := ShardCollector{
		DefaultCollector{
			Integration: i,
			Session:     test.MockSession{},
		},
		"testID",
		"testHost",
	}

	cc.CollectMetrics()

}

package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
)

func Test_ShardCollector_GetEntity(t *testing.T) {
	cc := ShardCollector{
		ID:   "testCollector",
		Host: "testCollector",
	}

	i, _ := integration.New("testIntegration", "testVersion")

	e, err := cc.GetEntity(i)
	if err != nil {
		t.Error(err)
	}

	if e.Metadata.Name != "testCollector" {
		t.Errorf("Expected entity name testCollector, got %s", e.Metadata.Name)
	}

	if e.Metadata.Namespace != "shard" {
		t.Errorf("Expected entity namespace shard, got %s", e.Metadata.Namespace)
	}

}

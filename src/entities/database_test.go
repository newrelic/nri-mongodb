package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_DatabaseCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")

	cc := DatabaseCollector{
		DefaultCollector{
			Integration: i,
			Session:     test.MockSession{},
		},
		"testCollector",
	}

	e, err := cc.GetEntity()
	if err != nil {
		t.Error(err)
	}

	if e.Metadata.Name != "testCollector" {
		t.Errorf("Expected entity name testCollector, got %s", e.Metadata.Name)
	}

	if e.Metadata.Namespace != "database" {
		t.Errorf("Expected entity namespace database, got %s", e.Metadata.Namespace)
	}

}

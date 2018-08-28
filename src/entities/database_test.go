package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
)

func Test_DatabaseCollector_GetEntity(t *testing.T) {
	cc := DatabaseCollector{
		Name: "testCollector",
	}

	i, _ := integration.New("testIntegration", "testVersion")

	e, err := cc.GetEntity(i)
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

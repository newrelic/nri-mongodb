package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

func Test_ConfigCollector_GetEntity(t *testing.T) {
	cc := ConfigCollector{
		HostCollector{ConnectionInfo: &connection.Info{
			Username: "testHost",
			Host:     "testCollector",
		},
		},
	}

	i, _ := integration.New("testIntegration", "testVersion")

	e, err := cc.GetEntity(i)
	if err != nil {
		t.Error(err)
	}

	if e.Metadata.Name != "testCollector" {
		t.Errorf("Expected entity name testCollector, got %s", e.Metadata.Name)
	}

	if e.Metadata.Namespace != "config" {
		t.Errorf("Expected entity namespace config, got %s", e.Metadata.Namespace)
	}

}

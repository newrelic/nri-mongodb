package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

func Test_MongodCollector_GetEntity(t *testing.T) {
	cc := MongodCollector{
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

	if e.Metadata.Namespace != "mongod" {
		t.Errorf("Expected entity namespace mongod, got %s", e.Metadata.Namespace)
	}

}

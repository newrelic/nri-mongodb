package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_MongodCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")

	cc := MongodCollector{
		HostCollector{
			DefaultCollector{
				Integration: i,
				Session:     test.MockSession{},
			},
			"testCollector",
		},
	}

	e, err := cc.GetEntity()
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

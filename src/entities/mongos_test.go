package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_MongosCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")

	cc := MongosCollector{
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

	if e.Metadata.Namespace != "mongos" {
		t.Errorf("Expected entity namespace mongos, got %s", e.Metadata.Namespace)
	}

}

func Test_MongosCollector_CollectMetrics(t *testing.T) {
	i, _ := integration.New("test", "0.0.1")
	cc := MongosCollector{
		HostCollector{
			DefaultCollector{
				Integration: i,
				Session:     test.MockSession{},
			},
			"testCollector",
		},
	}

	cc.CollectMetrics()

}

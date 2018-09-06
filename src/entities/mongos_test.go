package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_MongosCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")

	mc := mongosCollector{
		hostCollector{
			defaultCollector{
				"testCollector",
				i,
				test.MockSession{},
			},
		},
	}

	e, err := mc.GetEntity()
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
	mc := mongosCollector{
		hostCollector{
			defaultCollector{
				"testCollector",
				i,
				test.MockSession{},
			},
		},
	}

	mc.CollectMetrics()

}

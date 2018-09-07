package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func Test_ConfigCollector_GetEntity(t *testing.T) {
	i, _ := integration.New("testIntegration", "testVersion")

	cc := configCollector{
		hostCollector{
			defaultCollector{
				"testCollector",
				i,
				test.MockSession{},
			},
		},
	}

	e, err := cc.GetEntity()
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

func Test_ConfigCollector_CollectMetrics(t *testing.T) {
	i, _ := integration.New("test", "0.0.1")
	cc := configCollector{
		hostCollector{
			defaultCollector{
				"testCollector",
				i,
				test.MockSession{},
			},
		},
	}

	cc.CollectMetrics()

}

package entities

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func TestCollectServerStatus(t *testing.T) {

	i, _ := integration.New("test", "1")
	c := MongodCollector{
		HostCollector{
			DefaultCollector{
				Session:     test.MockSession{},
				Integration: i,
			},
			"testMongod",
		},
	}

	e, _ := c.GetEntity()
	ms := e.NewMetricSet("testmetricset")

	err := CollectServerStatus(c, ms)
	if err != nil {
		t.Error(err)
	}

}

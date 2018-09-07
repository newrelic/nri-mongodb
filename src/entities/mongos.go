package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// mongosCollector is a storage struct which contains all the information
// needed to collect metrics and inventory for a given mongos
type mongosCollector struct {
	hostCollector
}

// GetEntity creates or returns an entity for the mongos
func (c *mongosCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
		return i.Entity(c.name, "mongos")
	}

	return nil, errors.New("nil integration")
}

// CollectInventory collects inventory
func (c *mongosCollector) CollectInventory() {
	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to create config entity: %v", err)
		return
	}
	c.collectInventory(e)
}

// CollectMetrics sets all the metrics for the mongos
func (c *mongosCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to create entity: %v", err)
	}

	ms := e.NewMetricSet("MongosSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	if err := collectServerStatus(c, ms); err != nil {
		log.Error("Collect failed: %v", err)
	}
}

// GetMongoses returns an array of MongosCollectors which will be collected
func GetMongoses(session connection.Session, integration *integration.Integration) ([]Collector, error) {
	type MongosUnmarshaller []struct {
		ID string `bson:"_id"`
	}

	var mu MongosUnmarshaller
	c := session.DB("config").C("mongos")
	if err := c.FindAll(&mu); err != nil {
		return nil, err
	}

	mongoses := make([]Collector, 0, len(mu))
	for _, mongos := range mu {
		hostPort := extractHostPort(mongos.ID)
		mongosSession, err := session.New(hostPort.Host, hostPort.Port)
		if err != nil {
			log.Error("Failed to connect to mongos server %s: %v", mongos.ID, err)
			continue
		}

		mc := &mongosCollector{
			hostCollector{
				defaultCollector{
					hostPort.Host,
					integration,
					mongosSession,
				},
			},
		}

		mongoses = append(mongoses, mc)
	}

	return mongoses, nil
}

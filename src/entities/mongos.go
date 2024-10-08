package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/v3/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// mongosCollector is a storage struct which contains all the information
// needed to collect metrics and inventory for a given mongos
type mongosCollector struct {
	hostCollector
}

// GetEntity creates or returns an entity for the mongos
func (c *mongosCollector) GetEntity() (*integration.Entity, error) {
	if c.entity != nil {
		return c.entity, nil
	}
	if i := c.GetIntegration(); i != nil {
		ekey, err := c.GetSessionEntityKey()
		if err != nil {
			return nil, err
		}
		clusterNameIDAttr := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
		e, err := i.EntityReportedBy(ekey, c.name, "mo-mongos", clusterNameIDAttr)
		c.entity = e
		return e, err
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
	if logError(err, "Failed to create entity: %v") {
		return
	}

	ms := e.NewMetricSet("MongosSample",
		attribute.Attribute{Key: "displayName", Value: e.Metadata.Name},
		attribute.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
		attribute.Attribute{Key: "clusterName", Value: ClusterName},
	)

	logError(collectServerStatus(c, ms), "Collect failed: %v")
}

// GetMongoses returns an array of MongosCollectors which will be collected
func GetMongoses(session connection.Session, integration *integration.Integration) ([]Collector, error) {
	log.Debug("Determining which mongoses to collect")
	type MongosUnmarshaller []struct {
		ID string `bson:"_id" json:"_id"`
	}

	var mu MongosUnmarshaller
	if err := session.DB("config").C("mongos").FindAll(&mu); err != nil {
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
					hostPort.Host + ":" + hostPort.Port,
					integration,
					mongosSession,
					nil,
				},
			},
		}

		mongoses = append(mongoses, mc)
	}

	return mongoses, nil
}

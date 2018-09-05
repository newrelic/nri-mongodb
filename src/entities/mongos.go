package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// MongosCollector is a storage struct which contains all the information
// needed to collect metrics and inventory for a given mongos
type MongosCollector struct {
	HostCollector
}

// GetEntity creates or returns an entity for the mongos
func (c MongosCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
		return i.Entity(c.Name, "mongos")
	}

	return nil, errors.New("nil integration")
}

// CollectMetrics sets all the metrics for the mongos
func (c MongosCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to create entity: %v", err)
	}

	ms := e.NewMetricSet("MongosSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	if err := CollectServerStatus(c, ms); err != nil {
		log.Error("Collect failed: %v", err)
	}
}

// GetMongoses returns an array of MongosCollectors which will be collected
func GetMongoses(session connection.Session, integration *integration.Integration) ([]*MongosCollector, error) {
	type MongosUnmarshaller []struct {
		ID string `bson:"_id"`
	}

	var mu MongosUnmarshaller
	c := session.DB("config").C("mongos")
	if err := c.FindAll(&mu); err != nil {
		return nil, err
	}

	mongoses := make([]*MongosCollector, 0, len(mu))
	for _, mongos := range mu {
		hostPort := extractHostPort(mongos.ID)
		ci := connection.DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		session, err := ci.CreateSession()
		if err != nil {
			log.Error("Failed to connect to mongos server %s: %v", mongos.ID, err)
			continue
		}

		mc := &MongosCollector{
			HostCollector{
				DefaultCollector{
					Session:     session,
					Integration: integration,
				},
				ci.Host,
			},
		}

		mongoses = append(mongoses, mc)
	}

	return mongoses, nil
}

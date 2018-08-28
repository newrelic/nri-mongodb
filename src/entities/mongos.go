package entities

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// MongosCollector is a storage struct which contains all the information
// needed to collect metrics and inventory for a given mongos
type MongosCollector struct {
	HostCollector
}

// GetEntity creates or returns an entity for the mongos
func (c MongosCollector) GetEntity() (*integration.Entity, error) {
	return c.GetIntegration().Entity(c.Name, "mongos")
}

// CollectMetrics sets all the metrics for the mongos
func (c MongosCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to create entity: %v", err)
	}

	var ss metrics.ServerStatus
	if err := c.Session.DB("admin").Run(map[interface{}]interface{}{"serverStatus": 1}, &ss); err != nil {
		log.Error("Failed to collect serverStatus metrics for entity %s: %v", e.Metadata.Name, err)
	}
	ms := e.NewMetricSet("MongosSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	if err := ms.MarshalMetrics(ss); err != nil {
		log.Error("Failed to marshal metrics for entity %s: %v", e.Metadata.Name, err)
	}
}

// GetMongoses returns an array of MongosCollectors which will be collected
func GetMongoses(session connection.Session, integration *integration.Integration) ([]*MongosCollector, error) {
	type MongosUnmarshaller []struct {
		ID string `bson:"_id"`
	}

	var mu MongosUnmarshaller
	c := session.DB("config").C("mongos")
	if err := c.Find(map[interface{}]interface{}{}).All(&mu); err != nil {
		return nil, err
	}

	mongoses := make([]*MongosCollector, len(mu))
	for i, mongos := range mu {
		hostPort := extractHostPort(mongos.ID)
		ci := connection.DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		session, err := ci.CreateSession()
		if err != nil {
			return nil, err
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

		mongoses[i] = mc
	}

	return mongoses, nil
}

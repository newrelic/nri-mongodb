package entities

import (
	"fmt"

	"github.com/globalsign/mgo"
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
func (c MongosCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "mongos")
}

// CollectMetrics sets all the metrics for the mongos
func (c MongosCollector) CollectMetrics(e *integration.Entity) {
	session, err := c.ConnectionInfo.CreateSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", c.ConnectionInfo.Host, err)
		return
	}

	var ss metrics.ServerStatus
	if err := session.DB("admin").Run(map[interface{}]interface{}{"serverStatus": 1}, &ss); err != nil {
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
func GetMongoses(session *mgo.Session) ([]*MongosCollector, error) {
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

		mc := &MongosCollector{
			HostCollector{ConnectionInfo: ci},
		}

		mongoses[i] = mc
	}

	return mongoses, nil
}

package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// MongodCollector is a storage struct with all the information needed
// to collect metrics and inventory for a mongod
type MongodCollector struct {
	HostCollector
}

// GetEntity creates or returns an entity for the mongod
func (c MongodCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
		return i.Entity(c.Name, "mongod")
	}

	return nil, errors.New("nil integration")
}

// CollectMetrics sets all the metrics for a mongod
func (c MongodCollector) CollectMetrics() {
	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to get entity: %v", err)
		return
	}

	ms := e.NewMetricSet("MongodSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	isReplSet, err := CollectIsMaster(c, ms)
	if err != nil {
		log.Error("Collect failed: %v", err)
	}

	if isReplSet {
		if err := CollectReplSetMetrics(c, ms); err != nil {
			log.Error("Collect failed: %v", err)
		}
	}

	if err := CollectServerStatus(c, ms); err != nil {
		log.Error("Collect failed: %v", err)
	}

	if err := CollectTop(c); err != nil {
		log.Error("Collect failed: %v", err)
	}

}

// GetMongods returns an array of MongodCollectors to collect
func GetMongods(shard *ShardCollector, integration *integration.Integration) ([]*MongodCollector, error) {
	hostPorts, _ := parseReplicaSetString(shard.Host)

	mongodCollectors := make([]*MongodCollector, len(hostPorts))
	for i, hostPort := range hostPorts {
		ci := connection.DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		session, err := ci.CreateSession()
		if err != nil {
			return nil, err
		}

		newMongodCollector := &MongodCollector{
			HostCollector{
				DefaultCollector{
					Integration: integration,
					Session:     session,
				},
				ci.Host,
			},
		}
		mongodCollectors[i] = newMongodCollector
	}

	return mongodCollectors, nil
}

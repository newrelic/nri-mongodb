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
		if err := CollectReplGetStatus(c, e.Metadata.Name, ms); err != nil {
			log.Error("Collect failed: %v", err)
		}

		if err := CollectReplGetConfig(c, e.Metadata.Name, ms); err != nil {
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
func GetMongods(shardHostString string, integration *integration.Integration) ([]*MongodCollector, error) {
	hostPorts, _ := parseReplicaSetString(shardHostString)

	mongodCollectors := make([]*MongodCollector, 0, len(hostPorts))
	for _, hostPort := range hostPorts {
		ci := connection.DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		session, err := ci.CreateSession()
		if err != nil {
			log.Error("Failed to connected to mongod server %s: %v", ci.Host, err)
			continue
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
		mongodCollectors = append(mongodCollectors, newMongodCollector)
	}

	return mongodCollectors, nil
}

package entities

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// MongodCollector is a storage struct with all the information needed
// to collect metrics and inventory for a mongod
type MongodCollector struct {
	HostCollector
}

// GetEntity creates or returns an entity for the mongod
func (c MongodCollector) GetEntity() (*integration.Entity, error) {
	return c.GetIntegration().Entity(c.Name, "mongod")
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

	var isMaster metrics.IsMaster
	err = c.Session.DB("admin").Run(map[interface{}]interface{}{"isMaster": 1}, &isMaster)
	if err != nil {
		log.Error("failed") // TODO remove this when split into functions
	}

	if err := ms.MarshalMetrics(isMaster); err != nil {
		log.Error("Failed to marshal isMaster metrics for entity %s: %v", e.Metadata.Name, err)

	}

	if isMaster.SetName != nil {
		if err := collectReplSetMetrics(ms, c.Session); err != nil {
			log.Error("Failed to collect repl set metrics for entity %s: %v", e.Metadata.Name, err)
		}
	}

	// TODO split off into functions so they can return separately
	var ss metrics.ServerStatus
	if err := c.Session.DB("admin").Run(map[interface{}]interface{}{"serverStatus": 1}, &ss); err != nil {
		log.Error("Failed to collect serverStatus metrics for entity %s: %v", e.Metadata.Name, err)
	}

	if err := ms.MarshalMetrics(ss); err != nil {
		log.Error("Failed to marshal metrics for entity %s: %v", e.Metadata.Name, err)
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

func collectReplSetMetrics(ms *metric.Set, session connection.Session) error {

	var replSetStatus metrics.ReplSetGetStatus
	err := session.DB("admin").Run(map[interface{}]interface{}{"replSetGetStatus": 1}, &replSetStatus)
	if err != nil {
		return err
	}

	// TODO finish this

	return nil

}

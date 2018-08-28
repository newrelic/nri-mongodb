package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// ConfigCollector is a storage struct which holds all the
// necessary information to collect a config  server
type ConfigCollector struct {
	HostCollector
}

// GetEntity creates or returns an entity for the config server
func (c ConfigCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "config")
}

// CollectMetrics collects and sets metrics for a config server
func (c ConfigCollector) CollectMetrics(e *integration.Entity) {
	session, err := c.ConnectionInfo.CreateSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", c.ConnectionInfo.Host, err)
		return
	}

	ms := e.NewMetricSet("MongoConfigServerSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	var isMaster metrics.IsMaster
	err = session.DB("admin").Run(map[interface{}]interface{}{"isMaster": 1}, &isMaster)
	if err != nil {
		log.Error("failed to collect isMaster metrics for %s", e.Metadata.Name)
	}

	if err := ms.MarshalMetrics(isMaster); err != nil {
		log.Error("Failed to marshal isMaster metrics for %s: %v", e.Metadata.Name, err)
	}

	if isMaster.SetName != nil {
		if err := collectReplSetMetrics(ms, c.ConnectionInfo, session); err != nil {
			log.Error("Failed to collect repl set metrics for %s: %v", e.Metadata.Name, err)
		}
	}

	var ss metrics.ServerStatus
	if err := session.DB("admin").Run(map[interface{}]interface{}{"serverStatus": 1}, &ss); err != nil {
		log.Error("Failed to collect serverStatus metrics for %s: %v", e.Metadata.Name, err)
	}

	if err := ms.MarshalMetrics(ss); err != nil {
		log.Error("Failed to marshal metrics for %s: %v", e.Metadata.Name, err)
	}

}

// GetConfigServers returns a list of ConfigCollectors to collect
func GetConfigServers(session connection.Session) ([]*ConfigCollector, error) {
	type ConfigUnmarshaller struct {
		Map struct {
			Config string
		}
	}

	var cu ConfigUnmarshaller
	if err := session.DB("admin").Run("getShardMap", &cu); err != nil {
		return nil, err
	}

	configServersString := cu.Map.Config
	if configServersString == "" {
		return nil, errors.New("config hosts string not defined")
	}
	configHostPorts, _ := parseReplicaSetString(configServersString)

	configCollectors := make([]*ConfigCollector, len(configHostPorts))
	for i, configHostPort := range configHostPorts {
		ci := connection.DefaultConnectionInfo()
		ci.Host = configHostPort.Host
		ci.Port = configHostPort.Port

		cc := &ConfigCollector{
			HostCollector{ConnectionInfo: ci},
		}
		configCollectors[i] = cc
	}

	return configCollectors, nil
}

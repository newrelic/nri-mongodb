package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// ConfigCollector is a storage struct which holds all the
// necessary information to collect a config  server
type ConfigCollector struct {
	HostCollector
}

// GetEntity creates or returns an entity for the config server
func (c ConfigCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
		return i.Entity(c.Name, "config")
	}

	return nil, errors.New("nil integration")
}

// CollectMetrics collects and sets metrics for a config server
func (c ConfigCollector) CollectMetrics() {
	e, err := c.GetEntity()

	ms := e.NewMetricSet("MongoConfigServerSample",
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
}

// GetConfigServers returns a list of ConfigCollectors to collect
func GetConfigServers(session connection.Session, integration *integration.Integration) ([]*ConfigCollector, error) {
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

	var configCollectors []*ConfigCollector // Creation can fail, so can't pre-allocate
	for _, configHostPort := range configHostPorts {
		ci := connection.DefaultConnectionInfo()
		ci.Host = configHostPort.Host
		ci.Port = configHostPort.Port

		session, err := ci.CreateSession()
		if err != nil {
			log.Error("Failed to connect to config server %s", ci.Host)
			continue
		}

		cc := &ConfigCollector{
			HostCollector{
				DefaultCollector{
					Session:     session,
					Integration: integration,
				},
				ci.Host,
			},
		}
		configCollectors = append(configCollectors, cc)
	}

	return configCollectors, nil
}

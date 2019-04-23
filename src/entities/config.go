package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// configCollector is a storage struct which holds all the
// necessary information to collect a config  server
type configCollector struct {
	hostCollector
}

// GetEntity creates or returns an entity for the config server
func (c *configCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
    clusterNameIDAttr := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
		return i.EntityReportedBy(c.GetSessionEntityKey(), c.name, "mo-config", clusterNameIDAttr)
	}

	return nil, errors.New("nil integration")
}

// CollectInventory collects inventory
func (c *configCollector) CollectInventory() {
	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to create config entity: %v", err)
		return
	}
	c.collectInventory(e)
}

// CollectMetrics collects and sets metrics for a config server
func (c *configCollector) CollectMetrics() {
	e, err := c.GetEntity()
	if logError(err, "Failed to create config entity: %v") {
		return
	}

	ms := e.NewMetricSet("MongoConfigServerSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	isReplSet, err := collectIsMaster(c, ms)
	logError(err, "Collect is master failed: %v")

	if isReplSet {
		logError(collectReplGetStatus(c, e.Metadata.Name, ms), "Get ReplSet status failed: %v")

		logError(collectReplGetConfig(c, e.Metadata.Name, ms), "Get ReplSet config failed: %v")
	}

	logError(collectServerStatus(c, ms), "Collect server status failed: %v")
}

// GetConfigServers returns a list of ConfigCollectors to collect
func GetConfigServers(session connection.Session, integration *integration.Integration) ([]Collector, error) {
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

	configCollectors := make([]Collector, 0, len(configHostPorts))
	for _, configHostPort := range configHostPorts {
		configSession, err := session.New(configHostPort.Host, configHostPort.Port)
		if err != nil {
			log.Error("Failed to connect to config server %s: %v", configHostPort.Host, err)
			continue
		}

		cc := &configCollector{
			hostCollector{
				defaultCollector{
					configHostPort.Host + ":" + configHostPort.Port,
					integration,
					configSession,
				},
			},
		}
		configCollectors = append(configCollectors, cc)
	}

	return configCollectors, nil
}

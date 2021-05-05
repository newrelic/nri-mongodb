package entities

import (
	"errors"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// mongodCollector is a storage struct with all the information needed
// to collect metrics and inventory for a mongod
type mongodCollector struct {
	hostCollector
}

// GetEntity creates or returns an entity for the mongod
func (c *mongodCollector) GetEntity() (*integration.Entity, error) {
	if c.entity != nil {
		return c.entity, nil
	}
	if i := c.GetIntegration(); i != nil {
		ekey, err := c.GetSessionEntityKey()
		if err != nil {
			return nil, err
		}
		clusterNameIDAttr := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
		e, err := i.EntityReportedBy(ekey, c.name, "mo-mongod", clusterNameIDAttr)
		c.entity = e
		return e, err
	}

	return nil, errors.New("nil integration")
}

// CollectInventory collects inventory
func (c *mongodCollector) CollectInventory() {
	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to create mongod entity: %v", err)
		return
	}
	c.collectInventory(e)
}

// CollectMetrics sets all the metrics for a mongod
func (c *mongodCollector) CollectMetrics() {
	e, err := c.GetEntity()
	if logError(err, "Failed to create mongod entity: %v") {
		return
	}

	ms := e.NewMetricSet("MongodSample",
		attribute.Attribute{Key: "displayName", Value: e.Metadata.Name},
		attribute.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
		attribute.Attribute{Key: "clusterName", Value: ClusterName},
	)

	isReplSet, err := collectIsMaster(c, ms)
	logError(err, "Collect is master failed: %v")

	if isReplSet {
		logError(collectReplGetStatus(c, e.Metadata.Name, ms), "Get ReplSet status failed: %v")
		logError(collectReplGetConfig(c, e.Metadata.Name, ms), "Get ReplSet config failed: %v")
	}

	logError(collectServerStatus(c, ms), "Collect server status failed: %v")
	logError(collectTop(c), "Collect top failed: %v")
}

// GetStandaloneMongod creates a mongod from a session
func GetStandaloneMongod(session connection.Session, integration *integration.Integration) Collector {
	standaloneMongodCollector := &mongodCollector{
		hostCollector{
			defaultCollector{
				fmt.Sprintf("%s:%s", session.Info().Host, session.Info().Port),
				integration,
				session,
				nil,
			},
		},
	}

	return standaloneMongodCollector
}

// GetShardMongods attempts to connect to a member of the shard to retrieve shard configuration information
func GetShardMongods(session connection.Session, shardHostString string, integration *integration.Integration) ([]Collector, error) {
	hostPorts, _ := parseReplicaSetString(shardHostString)

	// For each host in the shard, attempt to get the connection info for all hosts in the shard
	for _, hostPort := range hostPorts {
		session, err := session.New(hostPort.Host, hostPort.Port)
		if err != nil {
			log.Warn("Failed to connect to mongod server %s:%s: %v", hostPort.Host, hostPort.Port, err)
			continue
		}

		mongods, err := GetReplSetMongods(session, integration)
		if err != nil {
			log.Warn("Failed to retrieve mongod collectors for replica set from %s: %s", hostPort.Host, err)
			continue
		}

		return mongods, nil

	}

	return nil, fmt.Errorf("failed to connect any mongods in shard %s", shardHostString)
}

// GetReplSetMongods attempts to connect to a member of a rreplica set to eplica set to
func GetReplSetMongods(session connection.Session, integration *integration.Integration) ([]Collector, error) {
	var replSetConfig metrics.ReplSetGetConfig
	if err := session.DB("admin").Run(Cmd{"replSetGetConfig": 1}, &replSetConfig); err != nil {
		return nil, fmt.Errorf("run replSetGetConfig failed: %s", err)
	}

	mongodCollectors := make([]Collector, 0, len(replSetConfig.Config.Members))
	for _, member := range replSetConfig.Config.Members {
		var host, port string
		hostPort := strings.Split(*member.Host, ":")
		if len(hostPort) == 2 {
			host = hostPort[0]
			port = hostPort[1]
		} else {
			host = hostPort[0]
			port = "27017"
		}
		mongodSession, err := session.New(host, port)
		if err != nil {
			log.Error("Failed to connected to mongod server %s. Skipping entity creation: %v", host, err)
			continue
		}

		newMongodCollector := &mongodCollector{
			hostCollector{
				defaultCollector{
					fmt.Sprintf("%s:%s", host, port),
					integration,
					mongodSession,
					nil,
				},
			},
		}
		mongodCollectors = append(mongodCollectors, newMongodCollector)
	}

	return mongodCollectors, nil
}

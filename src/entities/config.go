package entities

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

type ConfigCollector struct {
	HostCollector
}

func (c ConfigCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "config")
}

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
	err = session.Run(map[interface{}]interface{}{"isMaster": 1}, &isMaster)
	if err != nil {
		log.Error("failed to collect isMaster metrics for %s", e.Metadata.Name)
	}

	ms.MarshalMetrics(isMaster)

	if isMaster.SetName != "" {
		collectReplSetMetrics(ms, c.ConnectionInfo, session)
	}

	var ss metrics.ServerStatus
	session.Run(map[interface{}]interface{}{"serverStatus": 1}, &ss)

	ms.MarshalMetrics(ss)

}

func GetConfigServers(session *mgo.Session) ([]*ConfigCollector, error) {
	type ConfigUnmarshaller struct {
		Map struct {
			Config string
		}
	}

	var cu ConfigUnmarshaller
	session.Run("getShardMap", &cu)

	configServersString := cu.Map.Config
	if configServersString == "" {
		return nil, errors.New("config hosts string not defined")
	}
	configHostPorts, _ := parseReplicaSetString(configServersString)

	var configCollectors []*ConfigCollector
	for _, configHostPort := range configHostPorts {
		ci := connection.DefaultConnectionInfo()
		ci.Host = configHostPort.Host
		ci.Port = configHostPort.Port

		cc := &ConfigCollector{
			HostCollector{ConnectionInfo: ci},
		}
		configCollectors = append(configCollectors, cc)
	}

	return configCollectors, nil
}

package entities

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

type ConfigCollector struct {
	HostCollector
}

func (c ConfigCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "config")
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
	configHostPorts := extractHostsFromReplicaSetString(configServersString)

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

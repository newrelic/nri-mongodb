package main

import (
	"errors"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"strings"
)

var (
	EventType = "MongoDBSample"
)

type Collector interface {
	CollectMetrics(*integration.Entity)
	CollectInventory(*integration.Entity)
	GetEntity(*integration.Integration) (*integration.Entity, error)
}

type hostPort struct {
	Host string
	Port string
}

type DefaultCollector struct{}

func (d DefaultCollector) CollectMetrics(*integration.Entity) {
	return
}

func (d DefaultCollector) CollectInventory(*integration.Entity) {
	return
}

func (d DefaultCollector) GetEntity() string {
	return ""
}

type HostCollector struct {
	DefaultCollector
	ConnectionInfo *ConnectionInfo
}

func (c HostCollector) CollectInventory(*integration.Entity) {
	// TODO write inventory collection code
}

type MongosCollector struct {
	HostCollector
}

func (c MongosCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "mongos")
}

func (c MongosCollector) CollectMetrics(e *integration.Entity) {
	session, err := c.ConnectionInfo.createSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", c.ConnectionInfo.Host, err)
		return
	}

	var ss serverStatus
	session.Run(map[interface{}]interface{}{"serverStatus": 1}, &ss)
	ms := e.NewMetricSet(EventType,
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	ms.MarshalMetrics(ss)
}

type MongodCollector struct {
	HostCollector
}

func (c MongodCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "mongod")
}

type ConfigCollector struct {
	HostCollector
}

func (c ConfigCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "config")
}

func getMongoses() ([]*MongosCollector, error) {

	type MongosUnmarshaller []struct {
		ID string `bson:"_id"`
	}

	connectionInfo := DefaultConnectionInfo()
	session, err := connectionInfo.createSession()
	if err != nil {
		return nil, err
	}

	var mu MongosUnmarshaller
	c := session.DB("config").C("mongos")
	c.Find(map[interface{}]interface{}{}).All(&mu)

	var mongoses []*MongosCollector
	for _, mongos := range mu {
		hostPort := extractHostPort(mongos.ID)
		ci := DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		mc := &MongosCollector{
			HostCollector{ConnectionInfo: ci},
		}

		mongoses = append(mongoses, mc)
	}

	return mongoses, nil
}

func getConfigServers() ([]*ConfigCollector, error) {
	type ConfigUnmarshaller struct {
		Map struct {
			Config string
		}
	}

	connectionInfo := DefaultConnectionInfo()
	session, err := connectionInfo.createSession()
	if err != nil {
		return nil, err
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
		ci := DefaultConnectionInfo()
		ci.Host = configHostPort.Host
		ci.Port = configHostPort.Port

		cc := &ConfigCollector{
			HostCollector{ConnectionInfo: ci},
		}
		configCollectors = append(configCollectors, cc)
	}

	return configCollectors, nil
}

func extractHostPort(hostPortString string) hostPort {
	hostPortArray := strings.SplitN(hostPortString, ":", 2)
	if hostPortArray[1] == "" {
		hostPortArray[1] = args.Port
	}

	return hostPort{Host: hostPortArray[0], Port: hostPortArray[1]}
}

func extractHostsFromReplicaSetString(rsString string) []hostPort {
	if strings.Contains(rsString, "/") {
		rsString = strings.Split(rsString, "/")[1]
	}

	hostPortStrings := strings.Split(rsString, ",")
	var hostPorts []hostPort
	for _, hostPortString := range hostPortStrings {
		hostPorts = append(hostPorts, extractHostPort(hostPortString))
	}

	return hostPorts

}

package main

import (
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

func (c ConfigCollector) getEntity(i *integration.Integration) (*integration.Entity, error) {
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
		host, port := extractHostPort(mongos.ID)
		ci := DefaultConnectionInfo()
		ci.Host = host
		ci.Port = port

		mc := &MongosCollector{
			HostCollector{ConnectionInfo: ci},
		}

		mongoses = append(mongoses, mc)
	}

	return mongoses, nil
}

func extractHostPort(hostPort string) (string, string) {
	hostPortArray := strings.SplitN(hostPort, ":", 2)
	if hostPortArray[1] == "" {
		hostPortArray[1] = args.Port
	}

	return hostPortArray[0], hostPortArray[1]
}

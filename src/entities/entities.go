package entities

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// Collector is an interface which represents an entity.
// A Collector knows how to collect itself through the CollectMetrics
// and CollectInventory methods.
type Collector interface {
	CollectMetrics()
	CollectInventory()
	GetEntity() (*integration.Entity, error)
	GetIntegration() *integration.Integration
}

type hostPort struct {
	Host string
	Port string
}

// DefaultCollector is the most basic implementation of the
// Collector interface, and can be inherited to create a minimal
// running version which creates no metrics or inventory
type DefaultCollector struct {
	Session     connection.Session
	Integration *integration.Integration
}

// CollectMetrics collects no metrics
func (d DefaultCollector) CollectMetrics() {
	return
}

// CollectInventory collects no inventory
func (d DefaultCollector) CollectInventory() {
	return
}

// GetEntity returns a dummy entity
func (d DefaultCollector) GetEntity() (*integration.Entity, error) {
	return d.GetIntegration().Entity("", "")
}

// GetIntegration returns the integration associated with the collector
func (d DefaultCollector) GetIntegration() *integration.Integration {
	return d.Integration
}

func extractHostPort(hostPortString string) hostPort {
	hostPortArray := strings.SplitN(hostPortString, ":", 2)
	if len(hostPortArray) == 1 {
		return hostPort{Host: hostPortArray[0], Port: ""} // TODO use a better default port?
	}

	return hostPort{Host: hostPortArray[0], Port: hostPortArray[1]}
}

func parseReplicaSetString(rsString string) ([]hostPort, string) {

	rsName := ""
	if strings.Contains(rsString, "/") {
		split := strings.Split(rsString, "/")
		rsName = split[0]
		rsString = split[1]
	}

	hostPortStrings := strings.Split(rsString, ",")
	hostPorts := make([]hostPort, len(hostPortStrings))
	for i, hostPortString := range hostPortStrings {
		hostPorts[i] = extractHostPort(hostPortString)
	}

	return hostPorts, rsName

}

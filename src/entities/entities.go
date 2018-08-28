package entities

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
)

// Collector is an interface which represents an entity.
// A Collector knows how to collect itself through the CollectMetrics
// and CollectInventory methods.
type Collector interface {
	CollectMetrics(*integration.Entity)
	CollectInventory(*integration.Entity)
	GetEntity(*integration.Integration) (*integration.Entity, error)
}

type hostPort struct {
	Host string
	Port string
}

// DefaultCollector is the most basic implementation of the
// Collector interface, and can be inherited to create a minimal
// running version which creates no metrics or inventory
type DefaultCollector struct{}

// CollectMetrics collects no metrics
func (d DefaultCollector) CollectMetrics(e *integration.Entity) {
	return
}

// CollectInventory collects no inventory
func (d DefaultCollector) CollectInventory(e *integration.Entity) {
	return
}

// GetEntity returns a dummy entity
func (d DefaultCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity("defaultEntity", "entity")
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
	var hostPorts []hostPort
	for _, hostPortString := range hostPortStrings {
		hostPorts = append(hostPorts, extractHostPort(hostPortString))
	}

	return hostPorts, rsName

}

package entities

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"strings"
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

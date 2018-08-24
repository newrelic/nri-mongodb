package entities

import (
	"errors"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
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
		return hostPort{Host: hostPortArray[0], Port: args.Port}
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

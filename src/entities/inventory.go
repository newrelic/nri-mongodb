package entities

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

type HostCollector struct {
	DefaultCollector
	ConnectionInfo *connection.ConnectionInfo
}

func (c HostCollector) CollectInventory(*integration.Entity) {
	// TODO write inventory collection code
}

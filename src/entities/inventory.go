package entities

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// HostCollector is a base collector for any entity that represents a specific host
type HostCollector struct {
	DefaultCollector
	ConnectionInfo *connection.Info
}

// CollectInventory collects all the inventory for a given host
func (c HostCollector) CollectInventory(*integration.Entity) {
	// TODO write inventory collection code
}

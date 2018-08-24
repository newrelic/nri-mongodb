package entities

type HostCollector struct {
	DefaultCollector
	ConnectionInfo *ConnectionInfo
}

func (c HostCollector) CollectInventory(*integration.Entity) {
	// TODO write inventory collection code
}

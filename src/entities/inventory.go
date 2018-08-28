package entities

// HostCollector is a base collector for any entity that represents a specific host
type HostCollector struct {
	DefaultCollector
	Name string
}

// CollectInventory collects all the inventory for a given host
func (c HostCollector) CollectInventory() {
	// TODO write inventory collection code
}

package entities

type MongodCollector struct {
	HostCollector
}

func (c MongodCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "mongod")
}

func getMongods(shard *ShardCollector) ([]*MongodCollector, error) {
	hostPorts := extractHostsFromReplicaSetString(shard.Host)

	var mongodCollectors []*MongodCollector
	for _, hostPort := range hostPorts {
		ci := DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		newMongodCollector := &MongodCollector{
			HostCollector{ConnectionInfo: ci},
		}
		mongodCollectors = append(mongodCollectors, newMongodCollector)
	}

	return mongodCollectors, nil
}

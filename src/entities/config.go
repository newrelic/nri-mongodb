package entities

type ConfigCollector struct {
	HostCollector
}

func (c ConfigCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "config")
}

func getConfigServers() ([]*ConfigCollector, error) {
	type ConfigUnmarshaller struct {
		Map struct {
			Config string
		}
	}

	connectionInfo := DefaultConnectionInfo()
	session, err := connectionInfo.createSession()
	if err != nil {
		return nil, err
	}

	var cu ConfigUnmarshaller
	session.Run("getShardMap", &cu)

	configServersString := cu.Map.Config
	if configServersString == "" {
		return nil, errors.New("config hosts string not defined")
	}
	configHostPorts := extractHostsFromReplicaSetString(configServersString)

	var configCollectors []*ConfigCollector
	for _, configHostPort := range configHostPorts {
		ci := DefaultConnectionInfo()
		ci.Host = configHostPort.Host
		ci.Port = configHostPort.Port

		cc := &ConfigCollector{
			HostCollector{ConnectionInfo: ci},
		}
		configCollectors = append(configCollectors, cc)
	}

	return configCollectors, nil
}

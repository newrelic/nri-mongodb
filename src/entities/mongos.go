package entities

type MongosCollector struct {
	HostCollector
}

func (c MongosCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ConnectionInfo.Host, "mongos")
}

func (c MongosCollector) CollectMetrics(e *integration.Entity) {
	session, err := c.ConnectionInfo.createSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", c.ConnectionInfo.Host, err)
		return
	}

	var ss serverStatus
	session.Run(map[interface{}]interface{}{"serverStatus": 1}, &ss)
	ms := e.NewMetricSet("MongosSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	ms.MarshalMetrics(ss)
}

func GetMongoses(session *mgo.Session) ([]*MongosCollector, error) {
	type MongosUnmarshaller []struct {
		ID string `bson:"_id"`
	}

	var mu MongosUnmarshaller
	c := session.DB("config").C("mongos")
	c.Find(map[interface{}]interface{}{}).All(&mu)

	var mongoses []*MongosCollector
	for _, mongos := range mu {
		hostPort := extractHostPort(mongos.ID)
		ci := DefaultConnectionInfo()
		ci.Host = hostPort.Host
		ci.Port = hostPort.Port

		mc := &MongosCollector{
			HostCollector{ConnectionInfo: ci},
		}

		mongoses = append(mongoses, mc)
	}

	return mongoses, nil
}

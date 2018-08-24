package entities

type ShardCollector struct {
	DefaultCollector
	ID   string
	Host string
}

func (c ShardCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ID, "shard")
}

func (c ShardCollector) CollectMetrics(e *integration.Entity) {
	session, err := c.ConnectionInfo.createSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", c.ConnectionInfo.Host, err)
		return
	}
}

func getShards() ([]*ShardCollector, error) {
	type ShardUnmarshaller []struct {
		ID   string `bson:"_id"`
		Host string `bson:"host"`
	}

	connectionInfo := DefaultConnectionInfo()
	session, err := connectionInfo.createSession()
	if err != nil {
		return nil, err
	}

	var su ShardUnmarshaller
	c := session.DB("config").C("shards")
	c.Find(map[interface{}]interface{}{}).All(&su)

	var shards []*ShardCollector
	for _, shard := range su {
		mc := &ShardCollector{ID: shard.ID, Host: shard.Host}

		shards = append(shards, mc)
	}

	return shards, nil
}

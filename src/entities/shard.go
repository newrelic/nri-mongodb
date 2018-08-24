package entities

import (
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

type ShardCollector struct {
	DefaultCollector
	ID   string
	Host string
}

func (c ShardCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ID, "shard")
}

func (c ShardCollector) CollectMetrics(e *integration.Entity) {
	return
}

func GetShards(session *mgo.Session) ([]*ShardCollector, error) {
	type ShardUnmarshaller []struct {
		ID   string `bson:"_id"`
		Host string `bson:"host"`
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

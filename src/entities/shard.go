package entities

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// ShardCollector is a storage sturct which holds the necessary
// information to collect all the metrics and inventory for a specific shard
type ShardCollector struct {
	DefaultCollector
	ID   string
	Host string
}

// GetEntity creates or returns an entity for the shard
func (c ShardCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.ID, "shard")
}

// CollectMetrics sets all the metrics for the shard
func (c ShardCollector) CollectMetrics(e *integration.Entity) {

	ms := e.NewMetricSet("MongoShardSample",
		metric.Attribute{
			Key: "id", Value: c.ID,
		},
	)

	replSetHosts, replSetName := parseReplicaSetString(c.Host)
	if replSetName != "" {
		if err := ms.SetMetric("shard.isReplSet", true, metric.GAUGE); err != nil {
			log.Error("Failed to set metric shard.isReplSet for entity %s", e.Metadata.Name)
		}
		if err := ms.SetMetric("replset.name", replSetName, metric.ATTRIBUTE); err != nil {
			log.Error("Failed to set metric replset.name for entity %s", e.Metadata.Name)
		}

		connectionInfo := connection.DefaultConnectionInfo()
		connectionInfo.Host = replSetHosts[0].Host
		connectionInfo.Port = replSetHosts[0].Port

		// TODO finish this
		_, err := connectionInfo.CreateSession()
		if err != nil {
			log.Error("Failed to connect to %s: %v", connectionInfo.Host, err)
			return
		}
	} else {
		if err := ms.SetMetric("shard.isReplSet", false, metric.GAUGE); err != nil {
			log.Error("Failed to set metric shard.isReplSet for entity %s", e.Metadata.Name)
		}
	}

}

// GetShards creates an array of ShardCollectors
func GetShards(session connection.Session) ([]*ShardCollector, error) {
	type ShardUnmarshaller []struct {
		ID   string `bson:"_id"`
		Host string `bson:"host"`
	}

	var su ShardUnmarshaller
	c := session.DB("config").C("shards")
	if err := c.Find(map[interface{}]interface{}{}).All(&su); err != nil {
		return nil, err
	}

	shards := make([]*ShardCollector, len(su))
	for i, shard := range su {
		mc := &ShardCollector{
			ID:   shard.ID,
			Host: shard.Host,
		}

		shards[i] = mc
	}

	return shards, nil
}

package entities

import (
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
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

	ms := e.NewMetricSet("MongoShardSample",
		metric.Attribute{
			Key: "id", Value: c.ID,
		},
	)

	replSetHosts, replSetName := parseReplicaSetString(c.Host)
	if replSetName != "" {
		ms.SetMetric("shard.isReplSet", true, metric.GAUGE)
		ms.SetMetric("replset.name", replSetName, metric.ATTRIBUTE)

		connectionInfo := connection.DefaultConnectionInfo()
		connectionInfo.Host = replSetHosts[0].Host
		connectionInfo.Port = replSetHosts[0].Port

		_, err := connectionInfo.CreateSession()
		if err != nil {
			log.Error("Failed to connect to %s: %v", connectionInfo.Host, err)
			return
		}
	} else {
		ms.SetMetric("shard.isReplSet", false, metric.GAUGE)
	}

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
		mc := &ShardCollector{
			ID:   shard.ID,
			Host: shard.Host,
		}

		shards = append(shards, mc)
	}

	return shards, nil
}

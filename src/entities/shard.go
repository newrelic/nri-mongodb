package entities

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// GetShards creates an array of ShardCollectors
func GetShards(session connection.Session, integration *integration.Integration) ([]string, error) {
	type ShardUnmarshaller []struct {
		ID   string `bson:"_id"`
		Host string `bson:"host"`
	}

	var su ShardUnmarshaller
	if err := session.DB("config").C("shards").FindAll(&su); err != nil {
		return nil, err
	}

	shards := make([]string, len(su))
	for i, shard := range su {
		shards[i] = shard.Host
	}

	return shards, nil
}

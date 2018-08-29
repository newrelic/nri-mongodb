package entities

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// CollectionCollector is a storage struct which holds all the
// necessary information to collect a collection
type CollectionCollector struct {
	DefaultCollector
	Name string
	DB   string
}

// GetEntity creates or returns an entity for a collection
func (c CollectionCollector) GetEntity() (*integration.Entity, error) {
	return c.GetIntegration().Entity(c.Name, "collection")
}

// CollectMetrics collects and sets the metrics for a collection
func (c CollectionCollector) CollectMetrics() {
	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to get entity: %v")
		return
	}

	ms := e.NewMetricSet("MongoCollectionSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	if err := CollectCollStats(c, ms); err != nil {
		log.Error("Collect failed: %v", err)
		return
	}

}

// GetCollections returns a list of CollectionCollectors which each collect a collection
func GetCollections(dbName string, session connection.Session, integration *integration.Integration) ([]*CollectionCollector, error) {
	names, err := session.DB(dbName).CollectionNames()
	if err != nil {
		return nil, err
	}

	collections := make([]*CollectionCollector, len(names))
	for i, name := range names {
		newCollection := &CollectionCollector{
			DefaultCollector: DefaultCollector{
				Integration: integration,
				Session:     session,
			},
			Name: name,
			DB:   dbName,
		}

		collections[i] = newCollection
	}

	return collections, nil
}

package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/filter"
)

// collectionCollector is a storage struct which holds all the
// necessary information to collect a collection
type collectionCollector struct {
	defaultCollector
	db string
}

// GetEntity creates or returns an entity for a collection
func (c *collectionCollector) GetEntity() (*integration.Entity, error) {
  if c.entity != nil {
    return c.entity, nil
  }
	if i := c.GetIntegration(); i != nil {
    ekey, err := c.GetSessionEntityKey()
    if err != nil {
      return nil, err
    }
    clusterNameIDAttr := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
    databaseNameIDAttr := integration.IDAttribute{Key: "databaseName", Value: c.db}
    e, err := i.EntityReportedBy(ekey, c.name, "mo-collection", clusterNameIDAttr, databaseNameIDAttr)
    c.entity = e
    return e, err
	}

	return nil, errors.New("nil integration")
}

// CollectInventory no-op
func (c *collectionCollector) CollectInventory() {
}

// CollectMetrics collects and sets the metrics for a collection
func (c *collectionCollector) CollectMetrics() {
	e, err := c.GetEntity()
	if logError(err, "Failed to create collection entity: %v") {
		return
	}

	ms := e.NewMetricSet("MongoCollectionSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
		metric.Attribute{Key: "database", Value: c.db},
	)

	logError(collectCollStats(c, ms), "Collect failed: %v")
}

// GetCollections returns a list of CollectionCollectors which each collect a collection
func GetCollections(dbName string, session connection.Session, integration *integration.Integration, filter *filter.DatabaseFilter) ([]Collector, error) {
	names, err := session.DB(dbName).CollectionNames()
	if err != nil {
		return nil, err
	}

	collections := make([]Collector, 0, len(names))
	for _, name := range names {
		if filter == nil || filter.CheckFilter(dbName, name) {
			newCollection := &collectionCollector{
				defaultCollector{
					name,
					integration,
					session,
          nil,
				},
				dbName,
			}

			collections = append(collections, newCollection)
		}
	}

	return collections, nil
}

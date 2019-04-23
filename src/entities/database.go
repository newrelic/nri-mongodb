package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/filter"
)

// databaseCollector is a storage struct containing all the
// necessary information to collect a database
type databaseCollector struct {
	defaultCollector
}

// GetEntity creates or returns an entity for a database
func (c *databaseCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
    clusterNameIDAttr := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
		return i.EntityReportedBy(c.GetSessionEntityKey(), c.name, "mo-database", clusterNameIDAttr)
	}

	return nil, errors.New("nil integration")

}

// CollectInventory no-op
func (c *databaseCollector) CollectInventory() {
}

// CollectMetrics collects and sets all the database's metrics
func (c *databaseCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if logError(err, "Failed to create database entity: %v") {
		return
	}

	ms := e.NewMetricSet("MongoDatabaseSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	logError(collectDbStats(c, ms), "Collect failed: %v")
}

// GetDatabases returns a list of DatabaseCollectors which each collect a specific database
func GetDatabases(session connection.Session, integration *integration.Integration, filter *filter.DatabaseFilter) ([]Collector, error) {
	type DatabaseListUnmarshaller struct {
		Databases []struct {
			Name string `bson:"name"`
		} `bson:"databases"`
	}

	var unmarshalledDatabaseList DatabaseListUnmarshaller
	if err := session.DB("admin").Run(Cmd{"listDatabases": 1}, &unmarshalledDatabaseList); err != nil {
		return nil, err
	}

	databases := make([]Collector, 0, len(unmarshalledDatabaseList.Databases))
	for _, database := range unmarshalledDatabaseList.Databases {
		if filter == nil || filter.CheckFilter(database.Name, "") {
			newDatabase := &databaseCollector{
				defaultCollector{
					database.Name,
					integration,
					session,
				},
			}

			databases = append(databases, newDatabase)
		}
	}

	return databases, nil
}

package entities

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/arguments"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// databaseCollector is a storage struct containing all the
// necessary information to collect a database
type databaseCollector struct {
	defaultCollector
}

// GetEntity creates or returns an entity for a database
func (c *databaseCollector) GetEntity() (*integration.Entity, error) {
	return c.GetIntegration().Entity(c.name, "database")
}

// CollectInventory no-op
func (c *databaseCollector) CollectInventory() {
}

// CollectMetrics collects and sets all the database's metrics
func (c *databaseCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to get entity: %v", err)
	}

	ms := e.NewMetricSet("MongoDatabaseSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	if err := collectDbStats(c, ms); err != nil {
		log.Error("Collect failed: %s", err)
	}
}

// GetDatabases returns a list of DatabaseCollectors which each collect a specific database
func GetDatabases(session connection.Session, integration *integration.Integration, filter arguments.DatabaseFilter) ([]Collector, error) {
	type DatabaseListUnmarshaller struct {
		Databases []struct {
			Name string `bson:"name"`
		} `bson:"databases"`
	}

	var unmarshalledDatabaseList DatabaseListUnmarshaller
	if err := session.DB("admin").Run(cmd{"listDatabases": 1}, &unmarshalledDatabaseList); err != nil {
		return nil, err
	}

	//databases := make([]Collector, len(unmarshalledDatabaseList.Databases))
	databases := make([]Collector, 0)
	for _, database := range unmarshalledDatabaseList.Databases {
		if checkDatabaseFilter(database.Name, filter) {
			newDatabase := &databaseCollector{
				defaultCollector{
					database.Name,
					integration,
					session,
				},
			}

			//databases[i] = newDatabase
			databases = append(databases, newDatabase)
		} else {
			log.Info("Database '%s' did not match filter, will not collect", database.Name)
		}
	}

	return databases, nil
}

func checkDatabaseFilter(database string, filter arguments.DatabaseFilter) bool {
	// no filter, no whitelist
	if filter == nil {
		return true
	}

	for name := range filter {
		if name == database {
			return true
		}
	}
	return false
}
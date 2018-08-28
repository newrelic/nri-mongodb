package entities

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// DatabaseCollector is a storage struct containing all the
// necessary information to collect a database
type DatabaseCollector struct {
	DefaultCollector
	Name string
}

// GetEntity creates or returns an entity for a database
func (c DatabaseCollector) GetEntity() (*integration.Entity, error) {
	return c.GetIntegration().Entity(c.Name, "database")
}

// CollectMetrics collects and sets all the database's metrics
func (c DatabaseCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if err != nil {
		log.Error("Failed to get entity: %v", err)
	}

	var dbStats metrics.DbStats
	if err := c.Session.DB(c.Name).Run(map[interface{}]interface{}{"dbStats": 1}, &dbStats); err != nil {
		log.Error("Failed to collect dbStats metrics for entity %s: %v", e.Metadata.Name, err)
	}
	ms := e.NewMetricSet("MongoDatabaseSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	if err := ms.MarshalMetrics(dbStats); err != nil {
		log.Error("Failed to marshal dbStats metrics for entity %s: %v", e.Metadata.Name, err)
	}
}

// GetDatabases returns a list of DatabaseCollectors which each collect a specific database
func GetDatabases(session connection.Session, integration *integration.Integration) ([]*DatabaseCollector, error) {
	type DatabaseListUnmarshaller struct {
		Databases []struct {
			Name string `bson:"name"`
		} `bson:"databases"`
	}

	var unmarshalledDatabaseList DatabaseListUnmarshaller
	err := session.DB("admin").Run(map[interface{}]interface{}{"listDatabases": 1}, &unmarshalledDatabaseList)
	if err != nil {
		return nil, err
	}

	databases := make([]*DatabaseCollector, len(unmarshalledDatabaseList.Databases))
	for i, database := range unmarshalledDatabaseList.Databases {
		newDatabase := &DatabaseCollector{
			Name: database.Name,
		}

		databases[i] = newDatabase
	}

	return databases, nil
}

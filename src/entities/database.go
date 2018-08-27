package entities

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

type DatabaseCollector struct {
	DefaultCollector
	Name string
}

func (c DatabaseCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.Name, "database")
}

func (c DatabaseCollector) CollectMetrics(e *integration.Entity) {
	connectionInfo := connection.DefaultConnectionInfo()
	session, err := connectionInfo.CreateSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", connectionInfo.Host, err)
		return
	}

	var dbStats metrics.DbStats
	if err := session.DB(c.Name).Run(map[interface{}]interface{}{"dbStats": 1}, &dbStats); err != nil {
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

func GetDatabases(session *mgo.Session) ([]*DatabaseCollector, error) {
	type DatabaseListUnmarshaller struct {
		Databases []struct {
			Name string `bson:"name"`
		} `bson:"databases"`
	}

	var unmarshalledDatabaseList DatabaseListUnmarshaller
	err := session.Run(map[interface{}]interface{}{"listDatabases": 1}, &unmarshalledDatabaseList)
	if err != nil {
		return nil, err
	}

	var databases []*DatabaseCollector
	for _, database := range unmarshalledDatabaseList.Databases {
		newDatabase := &DatabaseCollector{
			Name: database.Name,
		}

		databases = append(databases, newDatabase)
	}

	return databases, nil
}

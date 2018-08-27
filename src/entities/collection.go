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

type CollectionCollector struct {
	DefaultCollector
	Name string
	DB   string
}

func (c CollectionCollector) GetEntity(i *integration.Integration) (*integration.Entity, error) {
	return i.Entity(c.Name, "collection")
}

func (c CollectionCollector) CollectMetrics(e *integration.Entity) {
	connectionInfo := connection.DefaultConnectionInfo()
	session, err := connectionInfo.CreateSession()
	if err != nil {
		log.Error("Failed to connect to %s: %v", connectionInfo.Host, err)
		return
	}

	var collStats metrics.CollStats
	session.DB(c.DB).Run(map[interface{}]interface{}{"collStats": c.Name}, &collStats)
	ms := e.NewMetricSet("MongoCollectionSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	ms.MarshalMetrics(collStats)
}

func GetCollections(dbName string, session *mgo.Session) ([]*CollectionCollector, error) {
	names, err := session.DB(dbName).CollectionNames()
	if err != nil {
		return nil, err
	}

	var collections []*CollectionCollector
	for _, name := range names {
		newCollection := &CollectionCollector{
			Name: name,
			DB:   dbName,
		}

		collections = append(collections, newCollection)
	}

	return collections, nil
}

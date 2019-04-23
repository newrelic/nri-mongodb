package entities

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// clusterCollector is a storage struct which contains all the information
// needed to collect metrics and inventory for the cluster
type clusterCollector struct {
	defaultCollector
}

// GetEntity creates or returns an entity for the mongos
func (c *clusterCollector) GetEntity() (*integration.Entity, error) {
	if i := c.GetIntegration(); i != nil {
		return i.EntityReportedBy(c.GetSessionEntityKey(), c.name, "mo-cluster")
	}

	return nil, errors.New("nil integration")
}

// CollectInventory collects inventory
func (c *clusterCollector) CollectInventory() {
}

// CollectMetrics sets all the metrics for the mongos
func (c *clusterCollector) CollectMetrics() {

	e, err := c.GetEntity()
	if logError(err, "Failed to create entity: %v") {
		return
	}

	ms := e.NewMetricSet("MongoSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
	)

	logError(collectNumDatabases(c, ms), "Collect failed: %v")
}

// GetClusters returns an array of MongosCollectors which will be collected
func GetClusters(session connection.Session, integration *integration.Integration) ([]Collector, error) {

	type MongosUnmarshaller []struct {
		ID string `bson:"_id" json:"_id"`
	}

	var mu MongosUnmarshaller
	if err := session.DB("config").C("mongos").FindAll(&mu); err != nil {
		return nil, err
	}

	clusters := make([]Collector, 0, 1)
	clusterName := ClusterName

	cluster := &clusterCollector{
		defaultCollector{
			clusterName,
			integration,
			session,
		},
	}

	clusters = append(clusters, cluster)
	return clusters, nil
}

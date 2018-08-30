package entities

import (
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// CollectServerStatus collects serverStatus metrics
func CollectServerStatus(c Collector, ms *metric.Set) error {

	// Retrieve the session for the collector
	session, err := c.GetSession()
	if err != nil {
		return fmt.Errorf("invalid session: %v", err)
	}

	// Collect and unmarshal the result
	var ss metrics.ServerStatus
	if err := session.DB("admin").Run(map[string]interface{}{"serverStatus": 1}, &ss); err != nil {
		return fmt.Errorf("run serverStatus failed: %v", err)
	}

	// Insert the metrics into the metric set
	if err := ms.MarshalMetrics(ss); err != nil {
		return fmt.Errorf("marshal metrics on serverStatus failed: %v", err)
	}

	return nil
}

// CollectIsMaster collects isMaster metrics. Returns a boolean which
// is true if the session is connected to a replica set
func CollectIsMaster(c Collector, ms *metric.Set) (bool, error) {

	// Retrieve the session for the collector
	session, err := c.GetSession()
	if err != nil {
		return false, fmt.Errorf("invalid session: %v", err)
	}

	// Collect and unmarshal the result
	var isMaster metrics.IsMaster
	if err := session.DB("admin").Run(map[string]interface{}{"isMaster": 1}, &isMaster); err != nil {
		return false, fmt.Errorf("run isMaster failed: %v", err)
	}

	// Insert the metrics into the metric set
	if err := ms.MarshalMetrics(isMaster); err != nil {
		return false, fmt.Errorf("marshal metrics on isMaster failed: %v", err)
	}

	// Return whether the node is part of a replica set and an error
	return isMaster.SetName != nil, nil

}

// CollectReplSetMetrics collects replica set metrics
func CollectReplSetMetrics(c Collector, ms *metric.Set) error {

	// Retrieve the session for the collector
	session, err := c.GetSession()
	if err != nil {
		return fmt.Errorf("invalid session: %v", err)
	}

	// Collect and unmarshal the metrics
	var replSetStatus metrics.ReplSetGetStatus
	if err := session.DB("admin").Run(map[string]interface{}{"replSetGetStatus": 1}, &replSetStatus); err != nil {
		return err
	}

	// TODO Finish this

	return nil

}

// CollectTop collects top metrics
func CollectTop(c Collector) error {
	session, err := c.GetSession()
	if err != nil {
		return fmt.Errorf("invalid session: %v", err)
	}

	e, err := c.GetEntity()
	if err != nil {
		return fmt.Errorf("invalid entity: %v", err)
	}

	var topMetrics metrics.Top
	if err := session.DB("admin").Run(map[string]interface{}{"top": 1}, &topMetrics); err != nil {
		return fmt.Errorf("run serverStatus failed: %v", err)
	}

	for key, collectionStats := range topMetrics.Totals {
		databaseName := strings.SplitN(key, ".", 2)[0]
		collectionName := strings.SplitN(key, ".", 2)[1]

		ms := e.NewMetricSet("MongodTopSample",
			metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
			metric.Attribute{Key: "entityName", Value: fmt.Sprintf("%s:%s", e.Metadata.Namespace, e.Metadata.Name)},
			metric.Attribute{Key: "database", Value: databaseName},
			metric.Attribute{Key: "collection", Value: collectionName},
		)

		if err := ms.MarshalMetrics(collectionStats); err != nil {
			return fmt.Errorf("marshal metrics on top failed: %v", err)
		}

	}

	return nil
}

// CollectCollStats collects collStats
func CollectCollStats(c CollectionCollector, ms *metric.Set) error {

	session, err := c.GetSession()
	if err != nil {
		return fmt.Errorf("invalid session: %v", err)
	}

	var collStats metrics.CollStats
	if err := session.DB(c.DB).Run(map[string]interface{}{"collStats": c.Name}, &collStats); err != nil {
		return fmt.Errorf("run collStats failed: %v", err)
	}

	if err := ms.MarshalMetrics(collStats); err != nil {
		return fmt.Errorf("marshal metrics on collStats failed: %v", err)
	}

	return nil
}

// CollectDbStats collects dbStats
func CollectDbStats(c DatabaseCollector, ms *metric.Set) error {
	var dbStats metrics.DbStats
	if err := c.Session.DB(c.Name).Run(map[string]interface{}{"dbStats": 1}, &dbStats); err != nil {
		return fmt.Errorf("run dbStats failed: %s", err)
	}

	if err := ms.MarshalMetrics(dbStats); err != nil {
		return fmt.Errorf("marshal metrics for dbStats failed: %s", err)
	}

	return nil
}

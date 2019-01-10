package main

import (
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/entities"
	"github.com/newrelic/nri-mongodb/src/filter"
)

// StartCollectorWorkerPool starts a pool of workers to handle collecting each entity
// and returns a channel of collectors which the workers read off of
func StartCollectorWorkerPool(numWorkers int, wg *sync.WaitGroup) chan entities.Collector {
	wg.Add(numWorkers)

	collectorChan := make(chan entities.Collector, 100)
	for j := 0; j < numWorkers; j++ {
		go collectorWorker(collectorChan, wg)
	}

	return collectorChan
}

// collectorWorker reads a collector from the collector chan, then asynchronously
// collects that object's inventory and metrics
func collectorWorker(collectorChan chan entities.Collector, wg *sync.WaitGroup) {
	defer wg.Done()

	// Loop until collectorChan is empty and closed
	for {
		collector, ok := <-collectorChan
		if !ok {
			return
		}

		// Create a waitGroup for collecting inventory and metrics
		var inventoryMetricsWg sync.WaitGroup

		if args.HasInventory() {
			inventoryMetricsWg.Add(1)
			go func() {
				defer inventoryMetricsWg.Done()
				collector.CollectInventory()
			}()
		}

		if args.HasMetrics() {
			inventoryMetricsWg.Add(1)
			go func() {
				defer inventoryMetricsWg.Done()
				collector.CollectMetrics()
			}()
		}

		// Wait until inventory and metrics are done collecting
		inventoryMetricsWg.Wait()
	}
}

// FeedWorkerPool feeds the workers with the collectors that contain the info needed to collect each entity
func FeedWorkerPool(session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer close(collectorChan)

	// Create a wait group for each of the get*Collectors calls
	getWg := new(sync.WaitGroup)

	isStandaloneInstance, err := entities.IsStandaloneInstance(session)
	if err != nil {
		log.Error("Failed to determine whether the monitored instance is a standalone instance or a sharded instance: %s", err.Error())
		return
	}

	if isStandaloneInstance {
		getWg.Add(1)
		go createStandaloneMongodCollector(getWg, session, collectorChan, integration)

		getWg.Add(1)
		go createDatabaseCollectors(getWg, session, collectorChan, integration)
	} else {
		getWg.Add(1)
		go createClusterCollectors(getWg, session, collectorChan, integration)

		getWg.Add(1)
		go createMongosCollectors(getWg, session, collectorChan, integration)

		getWg.Add(1)
		go createConfigCollectors(getWg, session, collectorChan, integration)

		getWg.Add(1)
		go createShardCollectors(getWg, session, collectorChan, integration)

		getWg.Add(1)
		go createDatabaseCollectors(getWg, session, collectorChan, integration)
	}

	getWg.Wait()
}

func createClusterCollectors(wg *sync.WaitGroup, session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	clusters, err := entities.GetClusters(session, integration)
	if err != nil {
		log.Error("Failed to collect list of clusters: %v", err)
	}
	for _, cluster := range clusters {
		collectorChan <- cluster
	}
}

func createMongosCollectors(wg *sync.WaitGroup, session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	mongoses, err := entities.GetMongoses(session, integration)
	if err != nil {
		log.Error("Failed to collect list of Mongos hosts: %v", err)
	}
	for _, mongos := range mongoses {
		collectorChan <- mongos
	}
}

func createConfigCollectors(wg *sync.WaitGroup, session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	configServers, err := entities.GetConfigServers(session, integration)
	if err != nil {
		log.Error("Failed to collect list of config servers: %v", err)
	}
	for _, configServer := range configServers {
		collectorChan <- configServer
	}
}

func createStandaloneMongodCollector(wg *sync.WaitGroup, session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	mongod := entities.GetStandaloneMongod(session, integration)
	collectorChan <- mongod

}

func createShardCollectors(wg *sync.WaitGroup, session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	shards, err := entities.GetShards(session, integration)
	if err != nil {
		log.Error("Failed to collect list of shards: %v")
	}
	for _, shard := range shards {
		// Create Mongod Collectors
		wg.Add(1)
		go func(localShard string) {
			defer wg.Done()

			mongods, err := entities.GetMongods(session, localShard, integration)
			if err != nil {
				log.Error("Failed to collect list of mongods for shard %s", localShard)
			}
			for _, mongod := range mongods {
				collectorChan <- mongod
			}
		}(shard)
	}
}

func createDatabaseCollectors(wg *sync.WaitGroup, session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	// this error is checked when arguments are validated
	databaseFilter, _ := filter.ParseFilters(args.Filters)
	databases, err := entities.GetDatabases(session, integration, databaseFilter)
	if err != nil {
		log.Error("Failed to collect list of databases: %v", err)
	}
	for _, database := range databases {
		collectorChan <- database

		// Create Collection Collectors
		wg.Add(1)
		go func(database entities.Collector) {
			defer wg.Done()
			collections, err := entities.GetCollections(database.GetName(), session, integration, databaseFilter)
			if err != nil {
				log.Error("Failed to collect list of collections for database %s: %v", database.GetName(), err)
			}

			for _, collection := range collections {
				collectorChan <- collection
			}
		}(database)
	}
}

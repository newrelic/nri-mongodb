package main

import (
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	args "github.com/newrelic/nri-mongodb/src/arguments"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/entities"
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

		if args.GlobalArgs.HasInventory() {
			inventoryMetricsWg.Add(1)
			go func() {
				defer inventoryMetricsWg.Done()
				collector.CollectInventory()
			}()
		}

		if args.GlobalArgs.HasMetrics() {
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

	// Create a wait group for each of the Get____ calls
	var getWg sync.WaitGroup

	// Create Mongos Collectors
	getWg.Add(1)
	go func() {
		defer getWg.Done()

		mongoses, err := entities.GetMongoses(session, integration)
		if err != nil {
			log.Error("Failed to collect list of Mongos hosts: %v", err)
		}
		for _, mongos := range mongoses {
			collectorChan <- mongos
		}
	}()

	// Create Config Server Collectors
	getWg.Add(1)
	go func() {
		defer getWg.Done()

		configServers, err := entities.GetConfigServers(session, integration)
		if err != nil {
			log.Error("Failed to collect list of config servers: %v", err)
		}
		for _, configServer := range configServers {
			collectorChan <- configServer
		}
	}()

	// Create Shard Collectors and their associated Mongod Collectors
	getWg.Add(1)
	go func() {
		defer getWg.Done()

		shards, err := entities.GetShards(session, integration)
		if err != nil {
			log.Error("Failed to collect list of shards: %v")
		}
		for _, shard := range shards {
			log.Info(shard.Host)
			collectorChan <- shard

			// Create Mongod Collectors
			getWg.Add(1)
			go func(localShard *entities.ShardCollector) {
				defer getWg.Done()

				mongods, err := entities.GetMongods(localShard.Host, integration)
				if err != nil {
					log.Error("Failed to collect list of mongods for shard %s", shard.ID)
				}
				for _, mongod := range mongods {
					collectorChan <- mongod
				}
			}(shard)
		}
	}()

	// Create Database Collectors and associated Collection Collectors
	getWg.Add(1)
	go func() {
		defer getWg.Done()

		databases, err := entities.GetDatabases(session, integration)
		if err != nil {
			log.Error("Failed to collect list of databases: %v", err)
		}
		for _, database := range databases {
			collectorChan <- database

			// Create Collection Collectors
			getWg.Add(1)
			go func(database *entities.DatabaseCollector) {
				defer getWg.Done()
				collections, err := entities.GetCollections(database.Name, session, integration)
				if err != nil {
					log.Error("Failed to collect list of collections for database %s: %v", database.Name, err)
				}

				for _, collection := range collections {
					collectorChan <- collection
				}
			}(database)
		}
	}()

	getWg.Wait()
}

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
func StartCollectorWorkerPool(numWorkers int, wg *sync.WaitGroup) chan entities.Collector {
	wg.Add(numWorkers)

	collectorChan := make(chan entities.Collector, 100)
	for j := 0; j < numWorkers; j++ {
		go collectorWorker(collectorChan, wg)
	}

	return collectorChan
}

func collectorWorker(collectorChan chan entities.Collector, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		collector, ok := <-collectorChan
		if !ok {
			return
		}

		if args.GlobalArgs.HasInventory() {
			collector.CollectInventory()
		}

		if args.GlobalArgs.HasMetrics() {
			collector.CollectMetrics()
		}
	}
}

// FeedWorkerPool feeds the workers with the collectors that contain the info needed to collect each entity
func FeedWorkerPool(session connection.Session, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer close(collectorChan)

	mongoses, err := entities.GetMongoses(session, integration)
	if err != nil {
		log.Error("Failed to collect list of Mongos hosts: %v", err)
	}
	for _, mongos := range mongoses {
		collectorChan <- mongos
	}

	configServers, err := entities.GetConfigServers(session, integration)
	if err != nil {
		log.Error("Failed to collect list of config servers: %v", err)
	}
	for _, configServer := range configServers {
		collectorChan <- configServer
	}

	shards, err := entities.GetShards(session, integration)
	if err != nil {
		log.Error("Failed to collect list of shards: %v")
	}
	for _, shard := range shards {
		collectorChan <- shard

		mongods, err := entities.GetMongods(shard, integration)
		if err != nil {
			log.Error("Failed to collect list of mongods for shard %s", shard.ID)
		}
		for _, mongod := range mongods {
			collectorChan <- mongod
		}
	}

	databases, err := entities.GetDatabases(session, integration)
	if err != nil {
		log.Error("Failed to collect list of databases: %v", err)
	}
	for _, database := range databases {
		collectorChan <- database

		collections, err := entities.GetCollections(database.Name, session, integration)
		if err != nil {
			log.Error("Failed to collect list of collections for database %s: %v", database.Name, err)
		}

		for _, collection := range collections {
			collectorChan <- collection
		}
	}
}

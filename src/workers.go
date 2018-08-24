package main

import (
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	args "github.com/newrelic/nri-mongodb/src/arguments"
	"github.com/newrelic/nri-mongodb/src/entities"
	"sync"
)

func startCollectorWorkerPool(numWorkers int, wg *sync.WaitGroup, i *integration.Integration) chan entities.Collector {
	wg.Add(numWorkers)

	collectorChan := make(chan entities.Collector, 100)
	for j := 0; j < numWorkers; j++ {
		go collectorWorker(collectorChan, wg, i)
	}

	return collectorChan
}

func collectorWorker(collectorChan chan entities.Collector, wg *sync.WaitGroup, i *integration.Integration) {
	defer wg.Done()

	for {
		collector, ok := <-collectorChan
		if !ok {
			return
		}

		entity, err := collector.GetEntity(i)
		if err != nil {
			log.Error("Failed to create entity for collector %+v: %v", collector, err)
		}

		if args.GlobalArgs.HasInventory() {
			collector.CollectInventory(entity)
		}

		if args.GlobalArgs.HasMetrics() {
			collector.CollectMetrics(entity)
		}
	}
}

func feedWorkerPool(session *mgo.Session, collectorChan chan entities.Collector) {
	defer close(collectorChan)

	mongoses, err := entities.GetMongoses(session)
	if err != nil {
		log.Error("Failed to collect list of Mongos hosts: %v", err)
	}
	for _, mongos := range mongoses {
		collectorChan <- mongos
	}

	configServers, err := entities.GetConfigServers(session)
	if err != nil {
		log.Error("Failed to collect list of config servers: %v", err)
	}
	for _, configServer := range configServers {
		collectorChan <- configServer
	}

	shards, err := entities.GetShards(session)
	if err != nil {
		log.Error("Failed to collect list of shards: %v")
	}
	for _, shard := range shards {
		collectorChan <- shard

		mongods, err := entities.GetMongods(shard)
		if err != nil {
			log.Error("Failed to collect list of mongods for shard %s", shard.ID)
		}
		for _, mongod := range mongods {
			collectorChan <- mongod
		}
	}

}

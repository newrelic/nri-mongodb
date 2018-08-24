package main

import (
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"sync"
)

func startCollectorWorkerPool(numWorkers int, wg *sync.WaitGroup, i *integration.Integration) chan Collector {
	wg.Add(numWorkers)

	collectorChan := make(chan Collector, 100)
	for j := 0; j < numWorkers; j++ {
		go collectorWorker(collectorChan, wg, i)
	}

	return collectorChan
}

func collectorWorker(collectorChan chan Collector, wg *sync.WaitGroup, i *integration.Integration) {
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

		if args.HasInventory() {
			collector.CollectInventory(entity)
		}

		if args.HasMetrics() {
			collector.CollectMetrics(entity)
		}
	}
}

func feedWorkerPool(session *mgo.Session, collectorChan chan Collector) {
	defer close(collectorChan)

	mongoses, err := GetMongoses(session)
	if err != nil {
		log.Error("Failed to collect list of Mongos hosts: %v", err)
	}
	for _, mongos := range mongoses {
		collectorChan <- mongos
	}

	configServers, err := GetConfigServers(session)
	if err != nil {
		log.Error("Failed to collect list of config servers: %v", err)
	}
	for _, configServer := range configServers {
		collectorChan <- configServer
	}

	shards, err := GetShards(session)
	if err != nil {
		log.Error("Failed to collect list of shards: %v")
	}
	for _, shard := range shards {
		collectorChan <- shard

		mongods, err := GetMongods(shard)
		if err != nil {
			log.Error("Failed to collect list of mongods for shard %s", shard.ID)
		}
		for _, mongod := range mongods {
			collectorChan <- mongod
		}
	}

}

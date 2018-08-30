package main

import (
	"os"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	args "github.com/newrelic/nri-mongodb/src/arguments"
	"github.com/newrelic/nri-mongodb/src/connection"
)

const (
	integrationName    = "com.newrelic.mongodb"
	integrationVersion = "0.1.0"
)

func main() {

	// Create the integration
	mongoIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args.GlobalArgs))
	if err != nil {
		log.Error("Failed to create integration: %v", err)
		os.Exit(1)
	}

	// Set verbose level
	log.SetupLogging(args.GlobalArgs.Verbose)

	// Validate arguments
	if err := args.GlobalArgs.Validate(); err != nil {
		log.Error("Invalid arguments: %v", err)
		os.Exit(1)
	}

	// Connect to Mongo
	connectionInfo := connection.DefaultConnectionInfo() // TODO only use args in the main package
	session, err := connectionInfo.CreateSession()

	// Start workers
	var wg sync.WaitGroup
	collectorChan := StartCollectorWorkerPool(100, &wg)

	// Feed the worker pool with entities to be collected
	go FeedWorkerPool(session, collectorChan, mongoIntegration)

	// Wait for workers to finish
	wg.Wait()

	// Publish the results
	if err = mongoIntegration.Publish(); err != nil {
		log.Error("Failed to publish integration: %v", err)
		os.Exit(1)
	}

}

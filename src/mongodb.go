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

	mongoIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args.GlobalArgs))
	if err != nil {
		log.Error("Failed to create integration: %v", err)
		os.Exit(1)
	}

	log.SetupLogging(args.GlobalArgs.Verbose)
	if err := args.GlobalArgs.Validate(); err != nil {
		log.Error("Invalid arguments: %v", err)
		os.Exit(1)
	}

	connectionInfo := connection.DefaultConnectionInfo()
	session, err := connectionInfo.CreateSession()

	var wg sync.WaitGroup
	collectorChan := StartCollectorWorkerPool(50, &wg)

	go FeedWorkerPool(session, collectorChan, mongoIntegration)

	wg.Wait()

	if err = mongoIntegration.Publish(); err != nil {
		log.Error("Failed to publish integration: %v", err)
		os.Exit(1)
	}

}

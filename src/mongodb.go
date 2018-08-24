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
	args.ValidateArguments(args.GlobalArgs)

	connectionInfo := connection.DefaultConnectionInfo()
	session, err := connectionInfo.CreateSession()

	var wg sync.WaitGroup
	// TODO change the worker pool size back to a higher number
	// after the concurrency panic bug is fixed
	collectorChan := startCollectorWorkerPool(1, &wg, mongoIntegration)

	go feedWorkerPool(session, collectorChan)

	wg.Wait()

	mongoIntegration.Publish()

}

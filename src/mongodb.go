package main

import (
	"os"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/arguments"
	"github.com/newrelic/nri-mongodb/src/connection"
)

const (
	integrationName    = "com.newrelic.mongodb"
	integrationVersion = "0.2.0"
)

var (
	args arguments.ArgumentList
)

func main() {
	// Create the integration
	mongoIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error("Failed to create integration")
		os.Exit(1)
	}

	// Set verbose level
	log.SetupLogging(args.Verbose)

	// Validate arguments
	if err := args.Validate(); err != nil {
		log.Error("Invalid arguments: %v", err)
		os.Exit(1)
	}

	// Connect to Mongo
	connectionInfo := connection.Info{
		AuthSource:            args.AuthSource,
		Host:                  args.Host,
		Password:              args.Password,
		Port:                  args.Port,
		Ssl:                   args.Ssl,
		SslCaCerts:            args.SslCaCerts,
		SslInsecureSkipVerify: args.SslInsecureSkipVerify,
		Username:              args.Username,
	}
	session, err := connectionInfo.CreateSession()
	if err != nil {
		log.Error("Failed to create session: %v", err)
		os.Exit(1)
	}

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

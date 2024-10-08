//go:generate goversioninfo
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
	"github.com/newrelic/nri-mongodb/src/arguments"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/entities"
)

const (
	integrationName = "com.newrelic.mongodb"
)

var (
	args               arguments.ArgumentList
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""
)

func main() {
	// Create the integration
	mongoIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error("Failed to create integration")
		os.Exit(1)
	}

	if args.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			strings.Title(strings.Replace(integrationName, "com.newrelic.", "", 1)),
			integrationVersion,
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(),
			gitCommit,
			buildDate)
		os.Exit(0)
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
		PEMKeyFile:            args.PEMKeyFile,
		Passphrase:            args.Passphrase,
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
	collectorChan := StartCollectorWorkerPool(args.ConcurrentCollections, &wg)

	// Set a global cluster name for identity attributes
	entities.ClusterName = args.MongodbClusterName

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

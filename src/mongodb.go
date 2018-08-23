package main

import (
	"os"
	"strconv"
	"sync"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Username              string `default:"" help:"Username for the MongoDB connection"`
	Password              string `default:"" help:"Password for the MongoDB connection"`
	Host                  string `default:"" help:"MongoDB host to connect to for monitoring"`
	Port                  string `default:"27017" help:"Port on which MongoDB is running"`
	AuthSource            string `default:"admin" help:"Database to authenticate against"`
	Ssl                   bool   `default:"false" help:"Enable SSL"`
	SslCertFile           string `default:"" help:"Path to the certificate file used to identify the local connection against MongoDB"`
	SslCaCerts            string `default:"" help:"Path to the ca_certs file"`
	SslInsecureSkipVerify bool   `default:"false" help:"Skip verification of the certificate sent by the host. This can make the connection susceptible to MITM attacks, and should only be used for testing."`
}

const (
	integrationName    = "com.newrelic.mongodb"
	integrationVersion = "0.1.0"
)

var (
	args argumentList
)

func main() {

	mongoIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error("Failed to create integration: %v", err)
		os.Exit(1)
	}

	log.SetupLogging(args.Verbose)
	validateArguments()

	var wg sync.WaitGroup
	collectorChan := startCollectorWorkerPool(10, &wg, mongoIntegration)

	connectionInfo := DefaultConnectionInfo()
	session, err := connectionInfo.createSession()

	go feedWorkerPool(session, collectorChan)

	wg.Wait()

	mongoIntegration.Publish()

}

func validateArguments() {
	if args.Username == "" {
		log.Error("Must provide a username argument")
		os.Exit(1)
	}

	if args.Password == "" {
		log.Error("Must provide a password argument")
		os.Exit(1)
	}

	if args.Host == "" {
		log.Error("Must provide a host argument")
		os.Exit(1)
	}

	if _, err := strconv.Atoi(args.Port); err != nil {
		log.Error("Invalid port %s", args.Port)
		os.Exit(1)
	}

	if args.SslInsecureSkipVerify {
		log.Warn("Using insecure SSL. This connection is susceptible to man in the middle attacks.")
	}

}

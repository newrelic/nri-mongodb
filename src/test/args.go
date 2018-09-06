package test

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-kafka/src/args"
)

// SetupTestArgs is a helper function for tests
func SetupTestArgs() {
	_, _ = integration.New("test", "0.0.1", integration.Args(&args.GlobalArgs))

}

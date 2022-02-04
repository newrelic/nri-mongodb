// +build integration

package integration

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/tests/integration/helpers"
	"github.com/newrelic/nri-mongodb/tests/integration/jsonschema"
	"github.com/stretchr/testify/assert"
)

var (
	iName = "mongodb"

	defaultContainer = "integration_nri-mongodb"

	defaultBinPath     = "/nri-mongodb"
	defaultHost        = "mongodb-sharded"
	defaultClusterName = "test-cluster"
	defaultUsername    = "root"
	defaultPassword    = "password123"

	// cli flags
	container = flag.String("container", defaultContainer, "container where the integration is installed")
	binPath   = flag.String("bin", defaultBinPath, "Integration binary path")

	username    = flag.String("username", defaultUsername, "Username for the MongoDB connection")
	password    = flag.String("password", defaultPassword, "Password for the MongoDB connection")
	host        = flag.String("host", defaultHost, "MongoDB host to connect to for monitoring")
	clusterName = flag.String("cluster_name", defaultClusterName, "A unique, user-defined name to identify the cluster")
)

// Returns the standard output, or fails testing if the command returned an error
func runIntegration(t *testing.T, envVars ...string) (string, string, error) {
	t.Helper()

	command := make([]string, 0)
	command = append(command, *binPath)

	var (
		hasEnvHost        bool
		hasEnvClusterName bool
		hasEnvUserName    bool
		hasEnvPassword    bool
	)

	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "HOST") {
			hasEnvHost = true
		}
		if strings.HasPrefix(envVar, "CLUSTER_NAME") {
			hasEnvClusterName = true
		}
		if strings.HasPrefix(envVar, "USERNAME") {
			hasEnvUserName = true
		}
		if strings.HasPrefix(envVar, "PASSWORD") {
			hasEnvPassword = true
		}
	}

	if !hasEnvHost && host != nil {
		command = append(command, "--host", *host)
	}
	if !hasEnvClusterName && clusterName != nil {
		command = append(command, "--mongodb_cluster_name", *clusterName)
	}
	if !hasEnvUserName && username != nil {
		command = append(command, "--username", *username)
	}
	if !hasEnvPassword && password != nil {
		command = append(command, "--password", *password)
	}

	stdout, stderr, err := helpers.ExecInContainer(*container, command, envVars...)

	if stderr != "" {
		log.Debug("Integration command Standard Error: ", stderr)
	}

	return stdout, stderr, err
}

func TestMain(m *testing.M) {
	flag.Parse()
	// Mongo cluster bootstrap wait time
	time.Sleep(20 * time.Second)
	result := m.Run()
	os.Exit(result)
}

func TestIntegrationConnectsToMongoV5(t *testing.T) {
	testName := helpers.GetTestName(t)
	_, stderr, err := runIntegration(t, "HOST=mongo5", "METRICS=true", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.Nil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")
}

func TestIntegrationMetrics(t *testing.T) {
	testName := helpers.GetTestName(t)
	stdout, stderr, err := runIntegration(t, "METRICS=true", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.NotNil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")
	schemaPath := filepath.Join("json-schema-files", "metrics-schema.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of integration doesn't have expected format.")
}
func TestIntegrationInventory(t *testing.T) {
	testName := helpers.GetTestName(t)
	stdout, stderr, err := runIntegration(t, "INVENTORY=true", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.NotNil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")
	schemaPath := filepath.Join("json-schema-files", "inventory-schema.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of integration doesn't have expected format.")
}
func TestIntegrationClusterName(t *testing.T) {
	testName := helpers.GetTestName(t)
	stdout, stderr, err := runIntegration(t, "METRICS=true", "MONGODB_CLUSTER_NAME=", "CLUSTER_NAME=test-cluster", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.NotNil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")
	schemaPath := filepath.Join("json-schema-files", "inventory-schema.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of integration doesn't have expected format.")
}

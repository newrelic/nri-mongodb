package entities

import (
	"errors"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/connection"
)

var (
	// ClusterName is an identifier for the cluster
	ClusterName string
)

// Cmd is an aliasi for map[string]interface{}
type Cmd map[string]interface{}

// Collector is an interface which represents an entity.
// A Collector knows how to collect itself through the CollectMetrics
// and CollectInventory methods.
type Collector interface {
	CollectMetrics()
	CollectInventory()
	GetName() string
	GetEntity() (*integration.Entity, error)
	GetIntegration() *integration.Integration
	GetSession() (connection.Session, error)
}

type hostPort struct {
	Host string
	Port string
}

func (d *defaultCollector) GetSessionEntityKey() (integration.EntityKey, error) {
	session, err := d.GetSession()
	if err != nil {
		return "", err
	}

	host := session.Info().Host
	port := session.Info().Port

	i := d.GetIntegration()
	clusterNameIDAttr := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
	var namespace string

	t, err := DetectDeploymentType(session)
	if err != nil {
		return "", err
	}

	if t == "standalone" || t == "replica_set" {
		namespace = "mo-mongod"
	} else {
		namespace = "mo-mongos"
	}
	e, err := i.Entity(fmt.Sprintf("%s:%s", host, port), namespace, clusterNameIDAttr)
	if err != nil {
		return "", err
	}
	key, err := e.Key()
	if err != nil {
		return "", err
	}
	return key, nil
}

// defaultCollector is the most basic implementation of the
// Collector interface, and can be inherited to create a minimal
// running version which creates no metrics or inventory
type defaultCollector struct {
	name        string
	integration *integration.Integration
	session     connection.Session
	entity      *integration.Entity
}

func (d *defaultCollector) GetName() string {
	return d.name
}

// GetIntegration returns the integration associated with the collector
func (d *defaultCollector) GetIntegration() *integration.Integration {
	return d.integration
}

// GetSession returns the session associated with the collector
func (d *defaultCollector) GetSession() (connection.Session, error) {
	if d.session != nil {
		return d.session, nil
	}

	return nil, errors.New("session is nil")
}

func logError(err error, format string, args ...interface{}) bool {
	if err != nil {
		if format == "" {
			log.Error("%v", err)
		} else {
			log.Error(format, append([]interface{}{err}, args...)...)
		}
		return true
	}
	return false
}

func extractHostPort(hostPortString string) hostPort {
	hostPortArray := strings.SplitN(hostPortString, ":", 2)
	if len(hostPortArray) == 1 || len(hostPortArray[1]) == 0 {
		return hostPort{Host: hostPortArray[0], Port: "27017"}
	}

	return hostPort{Host: hostPortArray[0], Port: hostPortArray[1]}
}

func parseReplicaSetString(rsString string) ([]hostPort, string) {

	rsName := ""
	if strings.Contains(rsString, "/") {
		split := strings.Split(rsString, "/")
		rsName = split[0]
		rsString = split[1]
	}

	hostPortStrings := strings.Split(rsString, ",")
	hostPorts := make([]hostPort, len(hostPortStrings))
	for i, hostPortString := range hostPortStrings {
		hostPorts[i] = extractHostPort(hostPortString)
	}

	return hostPorts, rsName

}

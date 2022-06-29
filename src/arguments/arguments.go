package arguments

import (
	"errors"
	"fmt"
	"strconv"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/filter"
)

// ArgumentList is a struct that defines the arguments for the integration
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Username              string `default:"" help:"Username for the MongoDB connection"`
	Password              string `default:"" help:"Password for the MongoDB connection"`
	Host                  string `default:"localhost" help:"MongoDB host to connect to for monitoring"`
	Port                  string `default:"27017" help:"Port on which MongoDB is running"`
	ClusterName           string `default:"" help:"(Deprecated in favor of MongodbClusterName)"`
	MongodbClusterName    string `default:"" help:"Cluster name to identify this Mongodb instance."`
	AuthSource            string `default:"admin" help:"Database to authenticate against"`
	Mechanism             string `default:"SCRAM-SHA-256" help:"Database authentication mechanism"`
	Ssl                   bool   `default:"false" help:"Enable SSL"`
	SslCaCerts            string `default:"" help:"Path to the ca_certs file"`
	PEMKeyFile            string `default:"" help:"PEM file contains Private Key and Client Certificate"`
	Passphrase            string `default:"" help:"Passphrase for decrypting Private Key"`
	SslInsecureSkipVerify bool   `default:"false" help:"Skip verification of the certificate sent by the host. This can make the connection susceptible to man-in-the-middle attacks, and should only be used for testing."`
	Filters               string `default:"" help:"JSON data defining database and collection filters."`
	ConcurrentCollections int    `default:"50" help:"The number of entities to collect metrics for concurrently. This is tunable to reduce CPU and memory requirements."`
	ShowVersion           bool   `default:"false" help:"Print build information and exit"`
}

// Validate validates an argument list and returns an error if something is wrong
func (args *ArgumentList) Validate() error {
	if args.Host == "" {
		return errors.New("must provide a host argument")
	}

	// ClusterName is being deprecated to avoid the collision with the nri-kubernetes integration.
	// For backward compatibility reasons the following fallback logic has been implemented to avoid breaking existant config.
	if args.MongodbClusterName == "" {
		if args.ClusterName == "" {
			return errors.New("Must supply a cluster name to identify this Mongodb cluster. Use MongodbClusterName config parameter")
		}
		args.MongodbClusterName = args.ClusterName
		log.Warn("Using the deprecated config ClusterName instead of MongodbClusterName")
	}

	if _, err := strconv.Atoi(args.Port); err != nil {
		return fmt.Errorf("invalid port %s", args.Port)
	}

	if args.ConcurrentCollections <= 0 {
		return fmt.Errorf("concurrent_collections must be greater than zero")
	}

	if args.SslInsecureSkipVerify {
		log.Warn("Using insecure SSL. This connection is susceptible to man in the middle attacks.")
	}

	if _, err := filter.ParseFilters(args.Filters); err != nil {
		return errors.New("invalid filter json argument")
	}

	return nil
}

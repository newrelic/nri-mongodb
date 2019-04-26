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
	ClusterName           string `default:"" help:"A unique, user-defined name to identify the cluster"`
	AuthSource            string `default:"admin" help:"Database to authenticate against"`
	Ssl                   bool   `default:"false" help:"Enable SSL"`
	SslCaCerts            string `default:"" help:"Path to the ca_certs file"`
	SslInsecureSkipVerify bool   `default:"false" help:"Skip verification of the certificate sent by the host. This can make the connection susceptible to man-in-the-middle attacks, and should only be used for testing."`
	Filters               string `default:"" help:"JSON data defining database and collection filters."`
}

// Validate validates an argument list and returns an error if something is wrong
func (args *ArgumentList) Validate() error {
	if args.Username == "" {
		return errors.New("must provide a username argument")
	}

	if args.Password == "" {
		return errors.New("must provide a password argument")
	}

	if args.Host == "" {
		return errors.New("must provide a host argument")
	}

	if args.ClusterName == "" {
		return errors.New("must provide a cluster_name argument")
	}

	if _, err := strconv.Atoi(args.Port); err != nil {
		return fmt.Errorf("invalid port %s", args.Port)
	}

	if args.SslInsecureSkipVerify {
		log.Warn("Using insecure SSL. This connection is susceptible to man in the middle attacks.")
	}

	if _, err := filter.ParseFilters(args.Filters); err != nil {
		return errors.New("invalid filter json argument")
	}

	return nil
}

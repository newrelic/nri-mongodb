package arguments

import (
	"errors"
	"fmt"
	"strconv"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
)

var (
	GlobalArgs ArgumentList
)

type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Username              string `default:"" help:"Username for the MongoDB connection"`
	Password              string `default:"" help:"Password for the MongoDB connection"`
	Host                  string `default:"localhost" help:"MongoDB host to connect to for monitoring"`
	Port                  string `default:"27017" help:"Port on which MongoDB is running"`
	AuthSource            string `default:"admin" help:"Database to authenticate against"`
	Ssl                   bool   `default:"false" help:"Enable SSL"`
	SslCertFile           string `default:"" help:"Path to the certificate file used to identify the local connection against MongoDB"`
	SslCaCerts            string `default:"" help:"Path to the ca_certs file"`
	SslInsecureSkipVerify bool   `default:"false" help:"Skip verification of the certificate sent by the host. This can make the connection susceptible to MITM attacks, and should only be used for testing."`
}

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

	if _, err := strconv.Atoi(args.Port); err != nil {
		return fmt.Errorf("invalid port %s", args.Port)
	}

	if args.SslInsecureSkipVerify {
		log.Warn("Using insecure SSL. This connection is susceptible to man in the middle attacks.")
	}

	return nil
}

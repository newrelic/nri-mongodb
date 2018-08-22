package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/globalsign/mgo"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Username    string `default:"" help:"Username for the MongoDB connection"`
	Password    string `default:"" help:"Password for the MongoDB connection"`
	Host        string `default:"" help:"MongoDB host to connect to for monitoring"`
	Port        string `default:"" help:"Port on which MongoDB is running"`
	AuthSource  string `default:"" help:"Database to authenticate against"`
	Ssl         bool   `default:"false" help:"Enable SSL"`
	SslCertFile string `default:"" help:"Path to the certificate file used to identify the local connection against MongoDB"`
	SslCaCerts  string `default:"" help:"Path to the ca_certs file"`
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

	// TODO validate arguments

	log.SetupLogging(args.Verbose)

	session, err := createSession()
	if err != nil {
		log.Error("Failed to create session: %v", err)
	}
	defer session.Close()

	var ss serverStatus
	err = session.Run(map[interface{}]interface{}{"serverStatus": 1}, &ss)
	if err != nil {
		log.Error("Failed to run command: %v", err)
	}
	fmt.Printf("%+v", ss)

	err = mongoIntegration.Publish()

}

func createSession() (*mgo.Session, error) {

	dialInfo := mgo.DialInfo{
		Addrs:    []string{args.Host},
		Username: args.Username,
		Password: args.Password,
	}

	if args.Ssl {
		tlsConfig := &tls.Config{}

		if args.SslCaCerts != "" {
			roots := x509.NewCertPool()

			ca, err := ioutil.ReadFile(args.SslCaCerts)
			if err != nil {
				log.Error("Failed to open crt file: %v", err)
			}

			roots.AppendCertsFromPEM(ca)

			tlsConfig.RootCAs = roots
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Error("Failed to dial Mongo instance: %v", err)
		os.Exit(1)
	}
	return session, err

}

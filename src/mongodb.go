package main

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Username   string `default:"" help:"The MongoDB connection user name"`
	Password   string `default:"" help:"The MongoDB connection password"`
	Hostname   string `default:"" help:"The MongoDB connection hostname"`
	AuthSource string `default:"" help:"Database to connect to"`
	Port       string `default:"" help:"Port to connect to"`
	Insecure   bool   `default:"true" help:"Indicates whether to skip the verification of the server certificate and hostname"`
	Ssl        bool   `default:"false" help:"Indicates whether SSL should be enabled"`
	// SslKeyFile  string `default:"" help:"Specifies the file containing the client certificate and private key used for authentication"`
	SslCertFile string `default:"" help:"Path to the certificate file used to identify the local connection against mongodb"`
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

	_, _ = integration.New(integrationName, integrationVersion, integration.Args(&args))

	client, err := configureClient()
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	nrdb := client.Database("newrelic")
	results, err := nrdb.RunCommand(nil, map[string]interface{}{"dbStats": 1})
	if err != nil {
		panic(err)
	}
	resultsBytes, err := bson.Marshal(results)
	if err != nil {
		panic(err)
	}

	type dbStats struct {
		Raw map[string]struct {
			Db          string
			Collections int
			Views       int
		}
	}

	var dbs dbStats
	err = bson.Unmarshal(resultsBytes, &dbs)
	if err != nil {
		panic(err)
	}

	fmt.Println(dbs)
}

func configureClient() (*mongo.Client, error) {
	opt := configureOptions()
	creds := configureCredentials()
	client, err := createClient(opt, creds)
	return client, err
}

func configureOptions() clientopt.SSLOpt {
	return clientopt.SSLOpt{
		Enabled:                      args.Ssl,
		ClientCertificateKeyFile:     args.SslCertFile,
		ClientCertificateKeyPassword: func() string { return args.Password },
		Insecure:                     args.Insecure,
		CaFile:                       args.SslCaCerts,
	}
}

func configureCredentials() clientopt.Credential {
	return clientopt.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    args.AuthSource,
		Username:      args.Username,
		Password:      args.Password,
	}
}

func createClient(sslopt clientopt.SSLOpt, creds clientopt.Credential) (*mongo.Client, error) {
	url := "mongodb://" + args.Hostname + ":" + args.Port + "/" + args.AuthSource
	client, err := mongo.NewClientWithOptions(
		url,
		clientopt.SSL(&sslopt),
		clientopt.Auth(creds),
	)
	if err != nil {
		log.Error("Error creating client")
		return nil, err
	}
	return client, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Username   string `default:"" help:"The MongoDB connection user name"`
	Password   string `default:"" help:"The MongoDB connection password"`
	Hostname   string `default:"" help:"The MongoDB connection hostname"`
	CaFile     string `default:"" help:"Specifies the file containing the certificate authority used for SSL connections"`
	AuthSource string `default:"" help:"Database to connect to"`
	Port       string `default:"" help:"Port to connect to"`
	KeyFile    string `default:"" help:"Specifies the file containing the client certificate and private key used for authentication"`
	EnableSSL  bool   `default:"false" help:"Indicates whether SSL should be enabled"`
	Insecure   bool   `default:"true" help:"Indicates whether to skip the verification of the server certificate and hostname"`
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

	var client *mongo.Client
	var err error
	if args.EnableSSL {
		client, err = configureSSLClient()
		if err != nil {
			println("handle error")
		}
	}

	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	nrdb := client.Database("admin")
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

func configureSSLClient() (*mongo.Client, error) {
	sslopt := configureSSLOptions()
	creds := configureSSLCredentials()
	client, err := createSSLClient(sslopt, creds)
	return client, err
}

func configureSSLOptions() clientopt.SSLOpt {
	return clientopt.SSLOpt{
		Enabled:  args.EnableSSL,
		Insecure: args.Insecure,
		CaFile:   args.CaFile,
		ClientCertificateKeyPassword: func() string { return args.Password },
		ClientCertificateKeyFile:     args.KeyFile,
	}
}

func configureSSLCredentials() clientopt.Credential {
	return clientopt.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    args.AuthSource,
		Username:      args.Username,
		Password:      args.Password,
	}
}

func createSSLClient(sslopt clientopt.SSLOpt, creds clientopt.Credential) (*mongo.Client, error) {
	println("mongodb://" + args.Hostname + ":" + args.Port + "/" + args.AuthSource)
	print(args.EnableSSL)
	url := "mongodb://" + args.Hostname + ":" + args.Port + "/" + args.AuthSource
	client, err := mongo.NewClientWithOptions(
		url,
		clientopt.SSL(&sslopt),
		clientopt.Auth(creds),
	)
	if err != nil {
		println(err)
		return nil, err
	}
	return client, nil
}

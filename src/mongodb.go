package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Username string
	Password string
}

const (
	integrationName    = "com.newrelic.mongodb"
	integrationVersion = "0.1.0"
)

var (
	args argumentList
)

func main() {

	sslopt := clientopt.SSLOpt{
		Enabled:  true,
		Insecure: false,
		CaFile:   "/Users/ccheek/bluemedora/blue_medora.crt",
	}

	creds := clientopt.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    "newrelic",
		Username:      "newrelic",
		Password:      "password",
	}

	client, err := mongo.NewClientWithOptions(
		"mongodb://mdb-rh7-rs1-r2.bluemedora.localnet:27017/newrelic",
		clientopt.SSL(&sslopt),
		clientopt.Auth(creds),
	)
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

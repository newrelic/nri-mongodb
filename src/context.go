package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoContext struct {
	connection *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

type Cmd = bson.D

func (this *MongoContext) Disconnect() error {
	if err := this.connection.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	return nil
}

func (this *MongoContext) DB(dbname string) *MongoContext {
	this.database = this.connection.Database(dbname)
	return this
}

func (this *MongoContext) C(colname string) *MongoContext {
	this.collection = this.database.Collection(colname)
	return this
}

func (this *MongoContext) Run(command Cmd) string {
	var result bson.M
	err := this.database.RunCommand(context.TODO(), command).Decode(&result)
	if err != nil {
		panic(err)
	}

	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(output)
}

func (this *MongoContext) ListDatabases() []string {
	result, err := this.connection.ListDatabaseNames(context.TODO(), Cmd{})
	if err != nil {
		panic(err)
	}
	return result
}

func (this *MongoContext) CollectionNames() []string {
	result, err := this.database.ListCollectionNames(context.TODO(), Cmd{})
	if err != nil {
		panic(err)
	}
	return result
}

func (this *MongoContext) Connect(uri string) *MongoContext {
	if uri == "" {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
	}

	var err error
	this.connection, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Initial state
	this.database = nil
	this.collection = nil

	return this
}

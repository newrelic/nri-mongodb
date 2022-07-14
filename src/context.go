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
	if this.connection == nil {
		return nil
	}
	if err := this.connection.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	this.connection = nil
	return nil
}

func (this *MongoContext) DB(dbname string) *MongoContext {
	if this.connection == nil {
		panic("DB needs a connection")
	}
	this.database = this.connection.Database(dbname)
	return this
}

func (this *MongoContext) C(colname string) *MongoContext {
	if this.database == nil {
		panic("C needs a database")
	}
	this.collection = this.database.Collection(colname)
	return this
}

func (this *MongoContext) FindAll(output interface{}) error {
	if this.collection == nil {
		panic("FindAll needs a collection")
	}
	cur, err := this.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return err
	}

	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), output); err != nil {
		return err
	}

	return nil
}

func (this *MongoContext) Run(command Cmd) string {
	if this.database == nil {
		panic("Run needs a database")
	}
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

func (this *MongoContext) RunUnmarshal(command Cmd, output interface{}) error {
	if this.database == nil {
		panic("RunUnmarshal needs a database")
	}
	var result bson.Raw
	err := this.database.RunCommand(context.TODO(), command).Decode(&result)
	if err != nil {
		return err
	}

	if err = bson.Unmarshal(result, output); err != nil {
		return err
	}
	return nil
}

func (this *MongoContext) ListDatabases() []string {
	if this.connection == nil {
		panic("ListDatabases needs a database")
	}
	result, err := this.connection.ListDatabaseNames(context.TODO(), Cmd{})
	if err != nil {
		panic(err)
	}
	return result
}

func (this *MongoContext) CollectionNames() []string {
	if this.database == nil {
		panic("CollectionNames needs a database")
	}
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

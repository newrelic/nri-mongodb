package connection

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session = *MongoConnection
type Cmd = bson.D

type MongoConnection struct {
	Host       string
	Port       string
	connection *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func (this *MongoConnection) Disconnect() error {
	if this.connection == nil {
		return nil
	}
	if err := this.connection.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	this.connection = nil
	return nil
}

func (this *MongoConnection) DB(dbname string) *MongoConnection {
	if this.connection == nil {
		panic("DB needs a connection")
	}
	this.database = this.connection.Database(dbname)
	return this
}

func (this *MongoConnection) C(colname string) *MongoConnection {
	if this.database == nil {
		panic("C needs a database")
	}
	this.collection = this.database.Collection(colname)
	return this
}

func (this *MongoConnection) FindAll(output interface{}) error {
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

func (this *MongoConnection) PipeAll(query []bson.M, output interface{}) error {
	if this.collection == nil {
		panic("FindAll needs a collection")
	}

	panic("Not implemented") // TODO: MongoDB Driver Port
}

func (this *MongoConnection) RunString(command Cmd) string {
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

func (this *MongoConnection) Run(command Cmd, output interface{}) error {
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

func (this *MongoConnection) ListDatabases() []string {
	if this.connection == nil {
		panic("ListDatabases needs a database")
	}
	result, err := this.connection.ListDatabaseNames(context.TODO(), Cmd{})
	if err != nil {
		panic(err)
	}
	return result
}

func (this *MongoConnection) CollectionNames() ([]string, error) {
	if this.database == nil {
		panic("CollectionNames needs a database")
	}
	result, err := this.database.ListCollectionNames(context.TODO(), Cmd{})
	if err != nil {
		panic(err)
	}
	return result, nil
}

func (this *MongoConnection) Connect(uri string) *MongoConnection {
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

func (this *MongoConnection) Info() *MongoConnection {
	return this
}

func (this *MongoConnection) New(host string, port string) (*MongoConnection, error) {
	// The new driver does not need clones
	return this, nil
}

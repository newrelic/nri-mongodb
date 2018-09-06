package test

import (
	"reflect"

	"github.com/globalsign/mgo/bson"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// MockSession is a mocked session
type MockSession struct{}

// DB returns a mocked DB
func (t MockSession) DB(name string) connection.DataLayer {
	return MockDB{}
}

// Close does nothing because this is a mock session
func (t MockSession) Close() {
	return
}

// New just returns itself because this is a mock session
func (t MockSession) New(host, port string) (connection.Session, error) {
	return t, nil
}

// MockDB is a mocked database
type MockDB struct{}

// C returns a mock collection
func (d MockDB) C(name string) connection.Collection {
	return MockCollection{}
}

// Run runs a command on a mock DB
func (d MockDB) Run(cmd interface{}, result interface{}) error {
	m := reflect.ValueOf(cmd)
	for _, key := range m.MapKeys() {
		switch key.String() {
		case "serverStatus":
			marshalled, _ := bson.Marshal(bson.M{
				"host":    "testhost",
				"version": "testversion",
				"process": "testprocess",
				"PID":     3538,
				"Uptime":  12345,
			})
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CollectionNames returns a mocked array of collection names
func (d MockDB) CollectionNames() ([]string, error) {
	return nil, nil
}

// MockCollection is a mock collection
type MockCollection struct{}

// FindAll runs a query on a mock collection
func (c MockCollection) FindAll(result interface{}) error {
	return nil
}

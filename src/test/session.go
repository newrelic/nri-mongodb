package test

import (
	"github.com/globalsign/mgo"
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

// MockDB is a mocked database
type MockDB struct{}

// C returns a mock collection
func (d MockDB) C(name string) connection.Collection {
	return MockCollection{}
}

// Run runs a command on a mock DB
func (d MockDB) Run(cmd interface{}, result interface{}) error {
	return nil

}

// CollectionNames returns a mocked array of collection names
func (d MockDB) CollectionNames() ([]string, error) {
	return nil, nil
}

// MockCollection is a mock collection
type MockCollection struct{}

// Find runs a query on a mock collection
func (c MockCollection) Find(query interface{}) *mgo.Query {
	// TODO figure out how to mock a Query
	return nil
}

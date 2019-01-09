package test

import (
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/stretchr/testify/mock"
)

// MockSession is a mockable connection.Session
type MockSession struct {
	mock.Mock
	databases map[string]*MockDatabase
}

// AssertExpectations asserts that everything specified with On and Return was
// in fact called as expected.  Calls may have occurred in any order.
func (s *MockSession) AssertExpectations(t mock.TestingT) bool {
	r := true
	r = r && s.Mock.AssertExpectations(t)
	for _, db := range s.databases {
		r = r && db.AssertExpectations(t)
	}
	return r
}

// MockDatabase returns, and adds if not present, a MockDatabase
func (s *MockSession) MockDatabase(name string, callCount int) *MockDatabase {
	db, ok := s.databases[name]
	if !ok {
		db = new(MockDatabase)
		s.On("DB", name).Return(db).Times(callCount)
		if s.databases == nil {
			s.databases = make(map[string]*MockDatabase)
		}
		s.databases[name] = db
	}
	return db
}

func (s *MockSession) Info() *connection.Info {
  return &connection.Info{}
}

// DB is mocked via setup
func (s *MockSession) DB(name string) connection.DataLayer {
	args := s.Called(name)
	return args.Get(0).(connection.DataLayer)
}

// New is mocked via setup
func (s *MockSession) New(host, port string) (connection.Session, error) {
	args := s.Called(host, port)
	session := args.Get(0)
	if session == nil {
		return nil, args.Error(1)
	}
	return session.(connection.Session), args.Error(1)
}

// Close is mocked via setup
func (s *MockSession) Close() {
	s.Called()
}

// MockDatabase is a mockable connection.DataLayer
type MockDatabase struct {
	mock.Mock
	collections map[string]*MockCollection
}

// AssertExpectations asserts that everything specified with On and Return was
// in fact called as expected.  Calls may have occurred in any order.
func (db *MockDatabase) AssertExpectations(t mock.TestingT) bool {
	r := true
	r = r && db.Mock.AssertExpectations(t)
	for _, coll := range db.collections {
		r = r && coll.AssertExpectations(t)
	}
	return r
}

// MockCollection returns, and adds if not present, a MockCollection
func (db *MockDatabase) MockCollection(name string, callCount int) *MockCollection {
	coll, ok := db.collections[name]
	if !ok {
		coll = new(MockCollection)
		db.On("C", name).Return(coll).Times(callCount)
		if db.collections == nil {
			db.collections = make(map[string]*MockCollection)
		}
		db.collections[name] = coll
	}
	return coll
}

// C is a mocked via setup
func (db *MockDatabase) C(name string) connection.Collection {
	args := db.Called(name)
	return args.Get(0).(connection.Collection)
}

// Run is mocked via setup
func (db *MockDatabase) Run(cmd interface{}, result interface{}) error {
	args := db.Called(cmd, result)
	return args.Error(0)
}

// CollectionNames is mocked via setup
func (db *MockDatabase) CollectionNames() ([]string, error) {
	args := db.Called()
	return args.Get(0).([]string), args.Error(1)
}

// MockCollection is a mockable connection.Collection
type MockCollection struct {
	mock.Mock
}

// FindAll is mocked via setup
func (c *MockCollection) FindAll(result interface{}) error {
	args := c.Called(result)
	return args.Error(0)
}

// PipeAll is mocked via setup
func (c *MockCollection) PipeAll(query, result interface{}) error {
	args := c.Called(query, result)
	return args.Error(0)
}

package connection

/*
 * Mockable Interfaces
 */

// Session is an interface that represents the
// minimum actions needed against a MongoDB Session.
type Session interface {
	DB(name string) DataLayer
	New(host, port string) (Session, error)
	Close()
}

// Collection is an interface that represents the minimum
// actions needed against a MongoDB collection
type Collection interface {
	FindAll(result interface{}) error
}

// DataLayer is an interface that represents the minimum
// actions needed against a MongoDB database
type DataLayer interface {
	C(name string) Collection
	Run(cmd interface{}, result interface{}) error
	CollectionNames() ([]string, error)
}

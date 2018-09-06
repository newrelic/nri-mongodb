package connection

import "github.com/globalsign/mgo"

// mongoDatabase is a struct that allows shadowing of mgo.Database functions for mocking
type mongoDatabase struct {
	*mgo.Database
}

// C is a function that shadows the C function of a mongo collection
func (d *mongoDatabase) C(name string) Collection {
	return &mongoCollection{d.Database.C(name)}
}

package connection

import "github.com/globalsign/mgo"

// mongoSession is a struct that allows shadowing of functions for mocking
type mongoSession struct {
	*mgo.Session
	info *Info
}

// DB shadows the mgo.Session DB function
func (s *mongoSession) DB(name string) DataLayer {
	return &mongoDatabase{s.Session.DB(name)}
}

// New creates a new session from the existing session. If the specified host and
// port are the same as the existing session, the existing session is returned.
func (s *mongoSession) New(host, port string) (Session, error) {
	if s.info.Host == host && s.info.Port == port {
		return s, nil
	}
	return s.info.clone(host, port).CreateSession()
}

package connection

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"time"

	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/arguments"
)

// SessionBuilder is a mockable interface that allows us to mock at the connection.Info level
type SessionBuilder interface {
	CreateSession() (Session, error)
}

// Info is a storage struct which holds all the
// information needed to connect to a Mongo host
type Info struct {
	Username              string
	Password              string
	AuthSource            string
	Host                  string
	Port                  string
	Ssl                   bool
	SslCaCerts            string
	SslInsecureSkipVerify bool
}

// Session is an interface that can be used to mock a MongoDB session
type Session interface {
	DB(name string) DataLayer
	Close()
}

// MongoSession is a struct that allows shadowing of functions for mocking
type MongoSession struct {
	*mgo.Session
}

// DB shadows the mgo.Session DB function
func (s *MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

// MongoDatabase is a struct that allows shadowing of mgo.Database functions for mocking
type MongoDatabase struct {
	*mgo.Database
}

// C is a function that shadows the C function of a mongo collection
func (d *MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// MongoCollection is a struct that allows shadowing of functions for mocking
type MongoCollection struct {
	*mgo.Collection
}

// DataLayer is an interface that can be used to mock a MongoDB database
type DataLayer interface {
	C(name string) Collection
	Run(cmd interface{}, result interface{}) error
	CollectionNames() ([]string, error)
}

// Collection is an interface that can be used to mock a MongoDB collection
type Collection interface {
	Find(query interface{}) *mgo.Query
}

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (Session, error) {

	host := c.Host
	if c.Port != "" {
		host += ":" + c.Port
	}
	dialInfo := mgo.DialInfo{
		Addrs:       []string{host},
		Username:    c.Username,
		Password:    c.Password,
		Source:      c.AuthSource,
		FailFast:    true,
		Timeout:     time.Duration(2) * time.Second,
		PoolTimeout: time.Duration(2) * time.Second,
		ReadTimeout: time.Duration(3) * time.Second,
	}

	if c.Ssl {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: c.SslInsecureSkipVerify,
		}

		if c.SslCaCerts != "" {
			roots := x509.NewCertPool()

			ca, err := ioutil.ReadFile(c.SslCaCerts)
			if err != nil {
				log.Error("Failed to open crt file: %v", err)
			}

			roots.AppendCertsFromPEM(ca)

			tlsConfig.RootCAs = roots
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, err
	}
	return &MongoSession{session}, nil
}

// DefaultConnectionInfo returns connection info constructed from the passed-in args
func DefaultConnectionInfo() *Info {
	connectionInfo := &Info{
		Username:              arguments.GlobalArgs.Username,
		Password:              arguments.GlobalArgs.Password,
		AuthSource:            arguments.GlobalArgs.AuthSource,
		Host:                  arguments.GlobalArgs.Host,
		Port:                  arguments.GlobalArgs.Port,
		Ssl:                   arguments.GlobalArgs.Ssl,
		SslCaCerts:            arguments.GlobalArgs.SslCaCerts,
		SslInsecureSkipVerify: arguments.GlobalArgs.SslInsecureSkipVerify,
	}

	return connectionInfo

}

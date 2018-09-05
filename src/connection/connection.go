package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/arguments"
)

/*
 * Mockable Interfaces
 */

// SessionBuilder is a mockable interface that allows us to mock at the connection.Info level
type SessionBuilder interface {
	CreateSession() (Session, error)
}

// Session is an interface that can be used to mock a MongoDB session
type Session interface {
	DB(name string) DataLayer
	Close()
}

// Collection is an interface that can be used to mock a MongoDB collection
type Collection interface {
	FindAll(result interface{}) error
}

// DataLayer is an interface that can be used to mock a MongoDB database
type DataLayer interface {
	C(name string) Collection
	Run(cmd interface{}, result interface{}) error
	CollectionNames() ([]string, error)
}

/*
 * Implementations of the mockable interfaces
 */

// MongoSession is a struct that allows shadowing of functions for mocking
type MongoSession struct {
	*mgo.Session
}

// DB shadows the mgo.Session DB function
func (s *MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{s.Session.DB(name)}
}

// MongoDatabase is a struct that allows shadowing of mgo.Database functions for mocking
type MongoDatabase struct {
	*mgo.Database
}

// C is a function that shadows the C function of a mongo collection
func (d *MongoDatabase) C(name string) Collection {
	return &MongoCollection{d.Database.C(name)}
}

// MongoCollection is a struct that allows shadowing of functions for mocking
type MongoCollection struct {
	col *mgo.Collection
}

// FindAll marshals all items in the collection into result
func (c *MongoCollection) FindAll(result interface{}) error {
	return c.col.Find(nil).All(result)
}

// Info is a storage struct which holds all the
// information needed to connect to a Mongo host.
// It implements the SessionBuilder interface
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

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (Session, error) {

	dialInfo := c.generateDialInfo()

	// TODO investigate this further. This should time out, but isn't.
	// The current manual timeout solution is dirty

	sessionChan := make(chan *mgo.Session)
	errChan := make(chan error)
	go func() {
		if session, err := mgo.DialWithInfo(dialInfo); err != nil {
			errChan <- err
		} else {
			sessionChan <- session
		}
	}()

	select {
	case session := <-sessionChan:
		return &MongoSession{session}, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(dialInfo.Timeout + (time.Duration(1) * time.Second)):
		return nil, fmt.Errorf("connection to %s timed out", dialInfo.Addrs[0])
	}

}

// generateDialInfo creates a dialInfo struct from a connection.Info struct
func (c *Info) generateDialInfo() *mgo.DialInfo {
	host := c.Host
	if c.Port != "" {
		host += ":" + c.Port
	}
	dialInfo := &mgo.DialInfo{
		Addrs:       []string{host},
		Username:    c.Username,
		Password:    c.Password,
		Source:      c.AuthSource,
		Direct:      true,
		FailFast:    true,
		Timeout:     time.Duration(10) * time.Second,
		PoolTimeout: time.Duration(10) * time.Second,
		ReadTimeout: time.Duration(10) * time.Second,
		ReadPreference: &mgo.ReadPreference{
			Mode: mgo.PrimaryPreferred,
		},
	}

	if c.Ssl {
		addSSL(dialInfo, c.SslInsecureSkipVerify, c.SslCaCerts)
	}

	return dialInfo
}

// addSSL adds SSL to a dialInfo struct
func addSSL(d *mgo.DialInfo, SslInsecureSkipVerify bool, SslCaCerts string) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: SslInsecureSkipVerify,
	}

	// If the user has defined a CA certificate file
	if SslCaCerts != "" {
		roots := x509.NewCertPool()

		ca, err := ioutil.ReadFile(SslCaCerts)
		if err != nil {
			log.Error("Failed to open SSL CA Certs file: %v", err)
		}

		roots.AppendCertsFromPEM(ca)

		tlsConfig.RootCAs = roots
	}

	// Use TLS to dial the server
	d.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
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

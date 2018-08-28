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

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (*mgo.Session, error) {

	// TODO figure out how port fits into here
	dialInfo := mgo.DialInfo{
		Addrs:       []string{c.Host},
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

	// TODO investigate this further. This should time out, but isn't.
	// The current manual timeout solution is dirty

	sessionChan := make(chan *mgo.Session)
	go func() {
		session, err := mgo.DialWithInfo(&dialInfo)
		if err != nil {
			log.Error("Failed to dial Mongo instance %s: %v", dialInfo.Addrs[0], err)
		}
		sessionChan <- session
	}()

	select {
	case session := <-sessionChan:
		return session, nil
	case <-time.After(time.Second * time.Duration(3)):
		return nil, fmt.Errorf("connection to %s timed out", dialInfo.Addrs[0])
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

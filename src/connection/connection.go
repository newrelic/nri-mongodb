package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type ConnectionInfo struct {
	Username              string
	Password              string
	AuthSource            string
	Host                  string
	Port                  string
	Ssl                   bool
	SslCaCerts            string
	SslInsecureSkipVerify bool
}

func DefaultConnectionInfo() *ConnectionInfo {
	newConnection := &ConnectionInfo{
		Username:              args.Username,
		Password:              args.Password,
		AuthSource:            args.AuthSource,
		Host:                  args.Host,
		Port:                  args.Port,
		Ssl:                   args.Ssl,
		SslCaCerts:            args.SslCaCerts,
		SslInsecureSkipVerify: args.SslInsecureSkipVerify,
	}

	return newConnection
}

func (c *ConnectionInfo) createSession() (*mgo.Session, error) {

	// TODO figure out how port fits into here
	dialInfo := mgo.DialInfo{
		Addrs:    []string{c.Host},
		Username: c.Username,
		Password: c.Password,
		Source:   c.AuthSource,
		FailFast: true,
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
		log.Error("Failed to dial Mongo instance %s: %v", dialInfo.Addrs[0], err)
		fmt.Printf("%+v\n", dialInfo)
		os.Exit(1)
	}
	return session, err

}

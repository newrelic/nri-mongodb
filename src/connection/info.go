package connection

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
)

// Info is a storage struct which holds all the
// information needed to connect to a Mongo host.
type Info struct {
	Username              string
	Password              string
	AuthSource            string
	Host                  string
	Port                  string
	Ssl                   bool
	SslCaCerts            string
	PEMKeyFile            string
	Passphrase            string
	SslInsecureSkipVerify bool
}

func (c *Info) clone(host, port string) *Info {
	if host == "" {
		host = c.Host
	}
	if port == "" {
		port = c.Port
	}
	info := &Info{
		Username:              c.Username,
		Password:              c.Password,
		AuthSource:            c.AuthSource,
		Host:                  host,
		Port:                  port,
		Ssl:                   c.Ssl,
		SslCaCerts:            c.SslCaCerts,
		PEMKeyFile:            c.PEMKeyFile,
		Passphrase:            c.Passphrase,
		SslInsecureSkipVerify: c.SslInsecureSkipVerify,
	}
	return info
}

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (Session, error) {

	dialInfo := c.generateDialInfo()

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
		return &mongoSession{Session: session, info: c}, nil
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
		addSSL(dialInfo, c.SslInsecureSkipVerify, c.SslCaCerts, c.PEMKeyFile, c.Passphrase)
	}

	return dialInfo
}

// addSSL adds SSL to a dialInfo struct
func addSSL(d *mgo.DialInfo, SslInsecureSkipVerify bool, SslCaCerts string, pemKeyFile string, passPhrase string) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: SslInsecureSkipVerify,
	}

	// If the user has defined a CA certificate file
	if SslCaCerts != "" {
		roots := x509.NewCertPool()

		if ca, err := ioutil.ReadFile(SslCaCerts); err != nil {
			log.Error("Failed to open SSL CA Certs file: %v", err)
		} else if !roots.AppendCertsFromPEM(ca) {
			log.Warn("No certificates were found in SSL CA certs file: %s", SslCaCerts)
		} else {
			tlsConfig.RootCAs = roots
		}
	}

	if pemKeyFile != "" {

		clientCert, err := parsePEMKeyFile(pemKeyFile, passPhrase)
		if err != nil {
			log.Error("Failed to open SSL PEM Key File: %v", err)
		} else {
			tlsConfig.Certificates = append([]tls.Certificate{}, clientCert)
		}
	}

	// Use TLS to dial the server
	d.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		if err != nil {
			log.Error("Failed to dial server: %s", err)
		}
		return conn, err
	}
}

func parsePEMKeyFile(pemKeyFile string, passphrase string) (tls.Certificate, error) {

	keyPemData, err := ioutil.ReadFile(pemKeyFile)
	if err != nil {
		return tls.Certificate{}, err
	}

	var pvtKeyBlock, clientCertBlock *pem.Block
	for {

		block, rest := pem.Decode(keyPemData)
		if block == nil {
			return tls.Certificate{}, errors.New("not a valid PEM file")
		}

		if strings.Contains(block.Type, "PRIVATE KEY") {
			pvtKeyBlock = block
		}

		if block.Type == "CERTIFICATE" {
			clientCertBlock = block
		}

		if len(rest) == 0 || (clientCertBlock != nil && pvtKeyBlock != nil) {
			break
		}

		keyPemData = rest
	}

	if clientCertBlock == nil {
		return tls.Certificate{}, errors.New("no Client Certificate found in PEM key file")
	}

	if pvtKeyBlock == nil {
		return tls.Certificate{}, errors.New("no Private Key found in PEM key file")
	}

	if x509.IsEncryptedPEMBlock(pvtKeyBlock) {
		decPemBytes, err := x509.DecryptPEMBlock(pvtKeyBlock, []byte(passphrase))
		if err != nil {
			return tls.Certificate{}, err
		}

		decBlock := pem.Block{}
		decBlock.Bytes = decPemBytes
		decBlock.Type = pvtKeyBlock.Type

		return tls.X509KeyPair(pem.EncodeToMemory(clientCertBlock), pem.EncodeToMemory(&decBlock))
	}

	return tls.X509KeyPair(pem.EncodeToMemory(clientCertBlock), pem.EncodeToMemory(pvtKeyBlock))
}

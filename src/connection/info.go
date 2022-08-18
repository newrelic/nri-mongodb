package connection

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/log"
)

// Info is a storage struct which holds all the
// information needed to connect to a Mongo host.
type Info struct {
	Username              string
	Password              string
	AuthSource            string
	Mechanism             string
	Host                  string
	Port                  string
	Ssl                   bool
	SslCaCerts            string
	PEMKeyFile            string
	Passphrase            string
	Atlas                 bool
	SslInsecureSkipVerify bool
	ConnectionString      string
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
		Mechanism:             c.Mechanism,
		Atlas:                 c.Atlas,
		Host:                  host,
		Port:                  port,
		Ssl:                   c.Ssl,
		SslCaCerts:            c.SslCaCerts,
		PEMKeyFile:            c.PEMKeyFile,
		Passphrase:            c.Passphrase,
		SslInsecureSkipVerify: c.SslInsecureSkipVerify,
		ConnectionString:      c.ConnectionString,
	}
	return info
}

func (c *Info) GetConnectionString() string {
	if c.ConnectionString == "" {
		c.ConnectionString = "mongodb"
		if c.Atlas {
			c.ConnectionString += "+srv"
		}
		c.ConnectionString += "://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/"
		if c.AuthSource != "" {
			c.ConnectionString += "?authSource=" + c.AuthSource
		}
		if c.Mechanism != "" {
			c.ConnectionString += "&authMechanism=" + c.Mechanism
		}
		if c.Ssl {
			if c.SslInsecureSkipVerify {
				c.ConnectionString += "&tlsInsecure=true"
			}
		}
	}

	return c.ConnectionString
}

func getSSL(SslInsecureSkipVerify bool, SslCaCerts string, pemKeyFile string, passPhrase string) *tls.Config {
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

	return tlsConfig
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

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (Session, error) {
	var conn MongoConnection
	conn.Host = c.Host
	conn.Port = c.Port
	var tlsConf *tls.Config = nil
	if c.Ssl {
		tlsConf = getSSL(c.SslInsecureSkipVerify, c.SslCaCerts, c.PEMKeyFile, c.Passphrase)
	}
	return conn.Connect(c.GetConnectionString(), tlsConf), nil
}

package connection

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"
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
		Mechanism:             c.Mechanism,
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

func (c *Info) ConnectionString() string {
	// TODO: Derive this from the Info struct
	return "mongodb://root:password123@localhost:27017"
}

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (Session, error) {
	var conn MongoConnection
	return conn.Connect(c.ConnectionString()), nil
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

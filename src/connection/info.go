package connection

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
		c.ConnectionString = "mongodb://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/"
		if c.AuthSource != "" {
			c.ConnectionString += "?authSource=" + c.AuthSource
		}
		if c.Mechanism != "" {
			c.ConnectionString += "&authMechanism=" + c.Mechanism
		}
		if c.Ssl {
			c.ConnectionString += "&ssl=True"
		}
		if c.PEMKeyFile != "" {
			c.ConnectionString += "&tlsCertificateKeyFile=" + c.PEMKeyFile
		}
		if c.SslInsecureSkipVerify {
			c.ConnectionString += "&tlsInsecure=true"
		}
	}
	return c.ConnectionString
}

// CreateSession uses the information in ConnectionInfo to return
// a session connected to a Mongo host
func (c *Info) CreateSession() (Session, error) {
	var conn MongoConnection
	conn.Host = c.Host
	conn.Port = c.Port
	return conn.Connect(c.GetConnectionString()), nil
}

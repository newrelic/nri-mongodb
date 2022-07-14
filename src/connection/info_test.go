package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo_clone(t *testing.T) {
	info := &Info{
		Username: "user",
		Password: "pwd",
		Host:     "host1",
		Port:     "1",
	}

	info2 := info.clone("", "")
	assert.Equal(t, info, info2, "Bad clone info2")

	info3 := info.clone("host3", "")
	info.Host = "host3"
	assert.Equal(t, info, info3, "Bad clone info3")

	info4 := info.clone("host4", "4")
	info.Host = "host4"
	info.Port = "4"
	assert.Equal(t, info, info4, "Bad clone info4")

	info5 := info.clone("", "5")
	info.Port = "5"
	assert.Equal(t, info, info5, "Bad clone info5")
}

// func TestInfo_CreateSession(t *testing.T) {
// 	info := &Info{
// 		Username:              "",
// 		Password:              "",
// 		Host:                  "localhost",
// 		Port:                  "27017",
// 		AuthSource:            "admin",
// 		Ssl:                   true,
// 		SslCaCerts:            "test",
// 		SslInsecureSkipVerify: false,
// 	}

// 	_, err := info.CreateSession()
// 	assert.Error(t, err, "Expected connection to fail")
// }

// func Test_addSSL(t *testing.T) {
// 	dialInfo := &mgo.DialInfo{
// 		Addrs:       []string{"localhost"},
// 		Username:    "",
// 		Password:    "",
// 		Source:      "admin",
// 		FailFast:    true,
// 		Timeout:     time.Duration(1) * time.Second,
// 		PoolTimeout: time.Duration(1) * time.Second,
// 		ReadTimeout: time.Duration(1) * time.Second,
// 	}

// 	addSSL(dialInfo, false, "", "", "")

// 	assert.NotNil(t, dialInfo.DialServer, "Nil dialServer")
// }

// func Test_addSSL_EmptyPEM(t *testing.T) {
// 	dialInfo := &mgo.DialInfo{
// 		Addrs:       []string{"localhost"},
// 		Username:    "",
// 		Password:    "",
// 		Source:      "admin",
// 		FailFast:    true,
// 		Timeout:     time.Duration(1) * time.Second,
// 		PoolTimeout: time.Duration(1) * time.Second,
// 		ReadTimeout: time.Duration(1) * time.Second,
// 	}

// 	addSSL(dialInfo, false, filepath.Join("testdata", "empty.pem"), "", "")

// 	assert.NotNil(t, dialInfo.DialServer, "Nil dialServer")
// }

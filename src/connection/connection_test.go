package connection

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/kr/pretty"
	"github.com/newrelic/nri-mongodb/src/arguments"
)

func TestDefaultConnectionInfo(t *testing.T) {
	arguments.GlobalArgs = arguments.ArgumentList{
		Host:       "localhost",
		Port:       "27017",
		AuthSource: "admin",
	}

	info := DefaultConnectionInfo()

	expected := &Info{
		Username:              "",
		Password:              "",
		Host:                  "localhost",
		Port:                  "27017",
		AuthSource:            "admin",
		Ssl:                   false,
		SslCaCerts:            "",
		SslInsecureSkipVerify: false,
	}

	if !reflect.DeepEqual(info, expected) {
		fmt.Println(pretty.Diff(info, expected))
		t.Error("Bad default arguments")
	}
}

func Test_generateDialInfo(t *testing.T) {
	arguments.GlobalArgs = arguments.ArgumentList{
		Host:       "localhost",
		Port:       "27017",
		AuthSource: "admin",
	}

	info := DefaultConnectionInfo()
	dialInfo := info.generateDialInfo()

	expectedDialInfo := &mgo.DialInfo{
		Addrs:       []string{"localhost:27017"},
		Username:    "",
		Password:    "",
		Source:      "admin",
		FailFast:    true,
		Timeout:     time.Duration(10) * time.Second,
		PoolTimeout: time.Duration(10) * time.Second,
		ReadTimeout: time.Duration(10) * time.Second,
	}

	if !reflect.DeepEqual(dialInfo, expectedDialInfo) {
		fmt.Println(pretty.Diff(dialInfo, expectedDialInfo))
		t.Error("Bad dial info")
	}

}

func Test_addSSL(t *testing.T) {
	dialInfo := &mgo.DialInfo{
		Addrs:       []string{"localhost"},
		Username:    "",
		Password:    "",
		Source:      "admin",
		FailFast:    true,
		Timeout:     time.Duration(10) * time.Second,
		PoolTimeout: time.Duration(10) * time.Second,
		ReadTimeout: time.Duration(10) * time.Second,
	}

	addSSL(dialInfo, false, "")

	if dialInfo.DialServer == nil {
		t.Error("Nil dialServer")
	}

}

func Test_Info_CreateSession(t *testing.T) {
	info := &Info{
		Username:              "",
		Password:              "",
		Host:                  "localhost",
		Port:                  "27017",
		AuthSource:            "admin",
		Ssl:                   true,
		SslCaCerts:            "test",
		SslInsecureSkipVerify: false,
	}

	_, err := info.CreateSession()
	if err == nil {
		t.Error("Expected connection to fail")
	}

}

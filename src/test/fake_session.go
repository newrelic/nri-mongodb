package test

import (
	"reflect"

	"github.com/globalsign/mgo/bson"
	"github.com/newrelic/infra-integrations-sdk/data/inventory"
	"github.com/newrelic/nri-mongodb/src/connection"
)

// FakeSession is a mocked session
type FakeSession struct{}

// DB returns a mocked DB
func (t FakeSession) DB(name string) connection.DataLayer {
	return FakeDB{}
}

// Close does nothing because this is a fake session
func (t FakeSession) Close() {
	return
}

// New just returns itself because this is a fake session
func (t FakeSession) New(host, port string) (connection.Session, error) {
	return t, nil
}

// FakeDB is a mocked database
type FakeDB struct{}

// C returns a fake collection
func (d FakeDB) C(name string) connection.Collection {
	return FakeCollection{}
}

// Run runs a command on a fake DB
func (d FakeDB) Run(cmd interface{}, result interface{}) error {
	m := reflect.ValueOf(cmd)
	for _, key := range m.MapKeys() {
		return unmarshalCommand(key.String(), result)
	}
	return nil
}

func unmarshalCommand(cmd string, result interface{}) error {
	switch cmd {
	case "serverStatus":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"asserts": map[string]interface{}{
				"regular":   100,
				"warning":   250,
				"msg":       600,
				"user":      3538,
				"rollovers": 12345,
			},
		})
		return bson.Unmarshal(marshalled, result)
	case "isMaster":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"setName":   "replica-set-name",
			"ismaster":  true,
			"secondary": true,
		})
		return bson.Unmarshal(marshalled, result)
	case "replSetGetStatus":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"name":     "mdb-rh7-rs1-a1.bluemedora.localnet:27017",
					"health":   1,
					"stateStr": "SECONDARY",
					"uptime":   758657,
				},
			},
		})
		return bson.Unmarshal(marshalled, result)
	case "replSetGetConfig":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"config": map[string]interface{}{
				"members": []map[string]interface{}{
					{
						"host":        "mdb-rh7-rs1-a1.bluemedora.localnet:27017",
						"arbiterOnly": 0.0,
						"hidden":      0.0,
						"priority":    10.0,
						"votes":       20.0,
					},
				},
			},
		})
		return bson.Unmarshal(marshalled, result)
	case "top":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"totals": map[string]interface{}{
				"records.users": map[string]interface{}{
					"total": map[string]interface{}{
						"time":  305277,
						"count": 2825,
					},
					"readLock": map[string]interface{}{
						"time":  305123,
						"count": 2893,
					},
					"writeLock": map[string]interface{}{
						"time":  13,
						"count": 1,
					},
				},
			},
		})
		return bson.Unmarshal(marshalled, result)
	case "collStats":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"size":       2157,
			"count":      3,
			"avgObjSize": 719,
			"capped":     false,
		})
		return bson.Unmarshal(marshalled, result)
	case "dbStats":
		marshalled, _ := bson.Marshal(map[string]interface{}{
			"objects":     5,
			"dataSize":    6,
			"storageSize": 7,
			"indexes":     4,
			"indexSize":   8,
		})
		return bson.Unmarshal(marshalled, result)
	case "getCmdLineOpts":
		return bson.UnmarshalJSON([]byte(`{
			"argv": [
				"/usr/bin/mongos",
				"-f",
				"/etc/mongodb.conf"
			],
			"parsed": {
				"config": "/etc/mongodb.conf",
				"net": {
					"bindIp": "0.0.0.0",
					"port": 27017
				},
				"systemLog": {
					"destination": "file",
					"logRotate": "rename"
				}
			},
			"ok": 1
		}`), result)
	case "getParameter":
		return bson.UnmarshalJSON([]byte(`{
			"one": 1,
			"two": ["one", "two"],
			"$three": "skipped",
			"ok": 1
		}`), result)
	}
	return nil
}

// CollectionNames returns a mocked array of collection names
func (d FakeDB) CollectionNames() ([]string, error) {
	return nil, nil
}

// FakeCollection is a fake collection
type FakeCollection struct{}

// FindAll runs a query on a fake collection
func (c FakeCollection) FindAll(result interface{}) error {
	return nil
}

// ExpectedInventory is what the inventory should look like using this fake session.
// It is a combination of "getCmdLineOpts" and "getParameter".
var ExpectedInventory = inventory.Items{
	"commandline/argv": {
		"0": "/usr/bin/mongos",
		"1": "-f",
		"2": "/etc/mongodb.conf",
	},
	"configuration/config": {
		"value": "/etc/mongodb.conf",
	},
	"configuration/net.bindIp": {
		"value": "0.0.0.0",
	},
	"configuration/net.port": {
		"value": float64(27017),
	},
	"configuration/systemLog.destination": {
		"value": "file",
	},
	"configuration/systemLog.logRotate": {
		"value": "rename",
	},
	"parameter/one": {
		"value": float64(1),
	},
	"parameter/two": {
		"value": "one, two",
	},
}

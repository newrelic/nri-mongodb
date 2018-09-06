package test

import (
	"fmt"

	"reflect"

	"github.com/globalsign/mgo/bson"
	"github.com/newrelic/nri-mongodb/src/connection"
	"github.com/newrelic/nri-mongodb/src/metrics"
)

// MockSession is a mocked session
type MockSession struct{}

// DB returns a mocked DB
func (t MockSession) DB(name string) connection.DataLayer {
	return MockDB{}
}

// Close does nothing because this is a mock session
func (t MockSession) Close() {
	return
}

// New just returns itself because this is a mock session
func (t MockSession) New(host, port string) (connection.Session, error) {
	return t, nil
}

// MockDB is a mocked database
type MockDB struct{}

// C returns a mock collection
func (d MockDB) C(name string) connection.Collection {
	return MockCollection{}
}

// Run runs a command on a mock DB
func (d MockDB) Run(cmd interface{}, result interface{}) error {
	m := reflect.ValueOf(cmd)
	for _, key := range m.MapKeys() {
		switch key.String() {
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

			err := bson.Unmarshal(marshalled, result)
			fmt.Printf("%v", result.(*metrics.ServerStatus).Asserts)
			if err != nil {
				return err
			}

		case "isMaster":
			marshalled, _ := bson.Marshal(map[string]interface{}{
				"ismaster":  true,
				"secondary": true,
			})
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}

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
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}

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
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}

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

			// marshalled, _ := bson.Marshal(map[string]interface{}{
			// 	"totals": map[string]interface{}{
			// 		"total": map[string]interface{}{
			// 			"time":  608,
			// 			"count": 1,
			// 		},
			// 	},
			// })
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}
		case "collStats":
			marshalled, _ := bson.Marshal(map[string]interface{}{
				"size":       2157,
				"count":      3,
				"avgObjSize": 719,
				"capped":     false,
			})
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}
		case "dbStats":
			marshalled, _ := bson.Marshal(map[string]interface{}{
				"objects":     5,
				"dataSize":    6,
				"storageSize": 7,
				"indexes":     4,
				"indexSize":   8,
			})
			err := bson.Unmarshal(marshalled, result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CollectionNames returns a mocked array of collection names
func (d MockDB) CollectionNames() ([]string, error) {
	return nil, nil
}

// MockCollection is a mock collection
type MockCollection struct{}

// FindAll runs a query on a mock collection
func (c MockCollection) FindAll(result interface{}) error {
	return nil
}

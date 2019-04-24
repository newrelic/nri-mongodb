package entities

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/mock"

	"github.com/newrelic/infra-integrations-sdk/data/inventory"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
	"github.com/stretchr/testify/assert"
)

func Test_hostCollector_collectInventory(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	e, _ := testIntegration.Entity("host", "namespace")
	mockSession := new(test.MockSession)
	mAdminDB := mockSession.MockDatabase("admin", 2)
	mAdminDB.On("Run", Cmd{"getCmdLineOpts": 1}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
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
			assert.NoError(t, err)
		}).
		Once()
	mAdminDB.On("Run", Cmd{"getParameter": "*"}, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(1)
			err := bson.UnmarshalJSON([]byte(`{
				"one": 1,
				"two": ["one", "two"],
				"$three": "skipped",
				"ok": 1
			}`), result)
			assert.NoError(t, err)
		}).
		Once()
	collector := &hostCollector{
		defaultCollector{
			"testHost",
			testIntegration,
			mockSession,
      nil,
		},
	}
	collector.collectInventory(e)

	mockSession.AssertExpectations(t)
	assert.Equal(t, test.ExpectedInventory, e.Inventory.Items())
}

func Test_hostCollector_collectInventory_Errors(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	e, _ := testIntegration.Entity("host", "namespace")
	mockSession := new(test.MockSession)
	mAdminDB := mockSession.MockDatabase("admin", 2)
	mAdminDB.On("Run", Cmd{"getCmdLineOpts": 1}, mock.Anything).
		Return(assert.AnError).
		Once()
	mAdminDB.On("Run", Cmd{"getParameter": "*"}, mock.Anything).
		Return(assert.AnError).
		Once()
	collector := &hostCollector{
		defaultCollector{
			"testHost",
			testIntegration,
			mockSession,
      nil,
		},
	}
	collector.collectInventory(e)
	expectedInventory := inventory.Items{}

	mockSession.AssertExpectations(t)
	assert.Equal(t, expectedInventory, e.Inventory.Items())
}

func Test_addInventoryArray(t *testing.T) {
	e := getTestEntity()
	arrayValue := []interface{}{
		"one",
		2,
		"three",
	}
	addInventoryArray(e, "", "first", arrayValue)
	addInventoryArray(e, "cat", "second", arrayValue, "p1")
	expectedInventory := inventory.Items{
		"first": {
			"value": "one, 2, three",
		},
		"cat/p1.second": {
			"value": "one, 2, three",
		},
	}
	assert.Equal(t, expectedInventory, e.Inventory.Items())
}

func Test_addInventoryMap(t *testing.T) {
	e := getTestEntity()
	mapValue := map[string]interface{}{
		"one": 1,
		"two": map[string]interface{}{
			"21": "2.1",
		},
		"three": []interface{}{
			"one",
			2,
			"three",
		},
	}
	addInventoryMap(e, "", mapValue, true)
	addInventoryMap(e, "cat", mapValue, true, "p1", "p2")
	expectedInventory := inventory.Items{
		"one": {
			"value": 1,
		},
		"two.21": {
			"value": "2.1",
		},
		"three": {
			"value": "one, 2, three",
		},
		"cat/p1.p2.one": {
			"value": 1,
		},
		"cat/p1.p2.two.21": {
			"value": "2.1",
		},
		"cat/p1.p2.three": {
			"value": "one, 2, three",
		},
	}

	assert.Equal(t, expectedInventory, e.Inventory.Items())
}

func Test_isInventoryKey(t *testing.T) {
	tests := []struct {
		key      string
		isRoot   bool
		expected bool
	}{
		{
			"ok",
			true,
			false,
		},
		{
			"ok",
			false,
			true,
		},
		{
			"operationTime",
			true,
			false,
		},
		{
			"operationTime",
			false,
			true,
		},
		{
			"$test",
			true,
			false,
		},
		{
			"$test",
			false,
			false,
		},
		{
			"test",
			true,
			true,
		},
		{
			"test",
			false,
			true,
		},
	}
	for _, tc := range tests {
		actual := isInventoryKey(tc.key, tc.isRoot)
		assert.Equal(t, tc.expected, actual)
	}
}

func Test_setInventoryItem(t *testing.T) {
	e := getTestEntity()
	setInventoryItem(e, "cat", "key", "field", "key-value1", "p1", "p2")
	setInventoryItem(e, "cat", "key", "field", "key-value2")
	setInventoryItem(e, "", "key", "field", "key-value3")
	expectedInventory := inventory.Items{
		"cat/p1.p2.key": {
			"field": "key-value1",
		},
		"cat/key": {
			"field": "key-value2",
		},
		"key": {
			"field": "key-value3",
		},
	}
	assert.Equal(t, expectedInventory, e.Inventory.Items())
}

func Test_setInventoryItem_Error(t *testing.T) {
	e := getTestEntity()
	tooLongKey := "lorem-ipsum-dolor-sit-amet-consectetur-adipiscing-elit-sed-do-eiusmod-tempor-incididunt-ut-labore-et-dolore-magna-aliqua-ut-enim-ad-minim-veniam-quis-nostrud-exercitation-ullamco-laboris-nisi-ut-aliquip-ex-ea-commodo-consequat-duis-aute-irure-dolor-in-reprehenderit-in-voluptate-velit-esse-cillum-dolore-eu-fugiat-nulla-pariatur-excepteur-sint-occaecat-cupidatat-non-proident-sunt-in-culpa-qui-officia-deserunt-mollit-anim-id-est-laborum"
	setInventoryItem(e, "", tooLongKey, "field", "value")
	assert.Empty(t, e.Inventory.Items())
}

func getTestEntity() *integration.Entity {
	i, _ := integration.New("test", "0.0.1")
	e, _ := i.Entity("host", "namespace")
	return e
}

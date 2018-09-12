package entities

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-mongodb/src/test"
)

func TestGetShards(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("config", 1).
		MockCollection("shards", 1).
		On("FindAll", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			result := args.Get(0)
			err := bson.UnmarshalJSON([]byte(`[
				{ "_id": "rs1", "host": "host1" },
				{ "_id": "rs2", "host": "host2" }
			]`), result)
			assert.NoError(t, err)
		}).
		Once()

	shards, err := GetShards(mockSession, testIntegration)
	expectedShards := []string{
		"host1",
		"host2",
	}

	mockSession.AssertExpectations(t)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedShards, shards)
}

func TestGetShards_Error(t *testing.T) {
	testIntegration, _ := integration.New("test", "0.0.1")
	mockSession := new(test.MockSession)
	mockSession.MockDatabase("config", 1).
		MockCollection("shards", 1).
		On("FindAll", mock.Anything).
		Return(assert.AnError).
		Once()

	shards, err := GetShards(mockSession, testIntegration)

	mockSession.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Nil(t, shards)
}

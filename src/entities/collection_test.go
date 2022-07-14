package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestCollectionCollector_GetEntity(t *testing.T) {
// 	cc := getTestCollectionCollector()

// 	e, err := cc.GetEntity()
// 	assert.NoError(t, err)
// 	assert.Equal(t, "testCollection", e.Metadata.Name)
// 	assert.Equal(t, "mo-collection", e.Metadata.Namespace)
// }

func TestCollectionCollector_GetEntity_Error(t *testing.T) {
	cc := getBadTestCollectionCollector()
	e, err := cc.GetEntity()
	assert.Nil(t, e)
	assert.Error(t, err)
}

// func TestCollectionCollector_CollectInventory(t *testing.T) {
// 	cc := getTestCollectionCollector()
// 	assert.NotPanics(t, func() {
// 		cc.CollectInventory()
// 	})
// }

// func TestCollectionCollector_CollectMetrics(t *testing.T) {
// 	cc := getTestCollectionCollector()

// 	assert.NotPanics(t, func() {
// 		cc.CollectMetrics()
// 	})
// }

func TestCollectionCollector_CollectMetrics_Error(t *testing.T) {
	cc := getBadTestCollectionCollector()

	assert.NotPanics(t, func() {
		cc.CollectMetrics()
	})
}

// FIXME: MongoDB Driver Port
// func TestGetCollections(t *testing.T) {
// 	testIntegration, _ := integration.New("test", "0.0.1")
// 	testFilter, _ := filter.ParseFilters("")
// 	expectedCollNames := []string{
// 		"CollectionOne",
// 		"CollectionTwo",
// 	}
// 	mockSession := new(test.MockSession)
// 	mockSession.MockDatabase("testDB", 1).
// 		On("CollectionNames").
// 		Return(expectedCollNames, nil).
// 		Once()

// 	collections, err := GetCollections("testDB", mockSession, testIntegration, testFilter)
// 	mockSession.AssertExpectations(t)
// 	assert.NoError(t, err)
// 	assert.Equal(t, len(expectedCollNames), len(collections))
// 	for i, coll := range collections {
// 		session, err := coll.GetSession()
// 		assert.NoError(t, err)
// 		assert.Equal(t, mockSession, session)
// 		assert.Equal(t, testIntegration, coll.GetIntegration())
// 		assert.Equal(t, expectedCollNames[i], coll.GetName())
// 	}
// }

// func TestGetCollections_Error(t *testing.T) {
// 	i, _ := integration.New("test", "0.0.1")
// 	testFilter, _ := filter.ParseFilters("")
// 	mockSession := new(test.MockSession)
// 	mockSession.MockDatabase("testDB", 1).
// 		On("CollectionNames").
// 		Return([]string{}, assert.AnError).
// 		Once()

// 	collections, err := GetCollections("testDB", mockSession, i, testFilter)
// 	mockSession.AssertExpectations(t)
// 	assert.Error(t, err)
// 	assert.Equal(t, assert.AnError, err)
// 	assert.Nil(t, collections)
// }

// func getTestCollectionCollector() *collectionCollector {
// 	i, _ := integration.New("test", "0.0.1")
// 	return &collectionCollector{
// 		defaultCollector{
// 			"testCollection",
// 			i,
// 			test.FakeSession{},
// 			nil,
// 		},
// 		"testDB",
// 	}
// }

func getBadTestCollectionCollector() *collectionCollector {
	return &collectionCollector{
		defaultCollector{
			"testCollection",
			nil,
			nil,
			nil,
		},
		"testDB",
	}
}

// Package filter contains a small helper that represents a whitelist of databases and collections to collect from
package filter

import (
	"encoding/json"
)

// DatabaseFilter represents a map of database and collection names to collect
type DatabaseFilter struct {
	CollectAll bool
	Filters    map[string][]string
}

// ParseFilters takes the given filter string json and parses it into a DatabaseFilter object
func ParseFilters(filterJSON string) (*DatabaseFilter, error) {
	// blank filter arg, no whitelist.
	if filterJSON == "" {
		return &DatabaseFilter{
			CollectAll: true,
			Filters: map[string][]string{},
		}, nil
	}

	filterMap := DatabaseFilter{
		CollectAll: false,
		Filters: map[string][]string{},
	}
	err := json.Unmarshal([]byte(filterJSON), &filterMap.Filters)
	if err != nil {
		return nil, err
	}
	return &filterMap, nil
}

// CheckFilter takes a database name and collection name and checks them against the whitelist
func (dbFilter *DatabaseFilter) CheckFilter(dbName, collectionName string) bool {
	// no filter, no whitelist
	if dbFilter.CollectAll {
		return true
	}

	for database, collections := range dbFilter.Filters {
		if database != dbName {
			continue
		}

		// if collectionName is empty, we're only filtering on the database
		if collectionName == "" {
			return true
		}
		
		return dbFilter.checkCollectionFilter(collections, collectionName)
	}

	// database was not included in whitelist, do not collect it
	return false
}

func (dbFilter *DatabaseFilter) checkCollectionFilter(collections []string, collectionName string) bool {
	// if a database is listed in the filter with a nil value, we want to collect all of its collections
	if collections == nil {
		return true
	}

	for _, collection := range collections {
		if collection == collectionName {
			return true
		}
	}

	return false
}
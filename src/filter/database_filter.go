// Package filter contains a small helper that represents a whitelist of databases and collections to collect from
package filter

import (
	"encoding/json"

	"github.com/newrelic/infra-integrations-sdk/log"
)

// DatabaseFilter represents a map of database and collection names to collect
type DatabaseFilter struct {
	CollectAll bool
	Filters    map[string]map[string]struct{}
}

// ParseFilters takes the given filter string json and parses it into a DatabaseFilter object
func ParseFilters(filterJSON string) (*DatabaseFilter, error) {
	// blank filter arg, no whitelist.
	if filterJSON == "" {
		return &DatabaseFilter{
			CollectAll: true,
			Filters:    nil,
		}, nil
	}

	filterMap := DatabaseFilter{
		CollectAll: false,
		Filters:    map[string]map[string]struct{}{},
	}

	jsonResult := make(map[string][]string)
	if err := json.Unmarshal([]byte(filterJSON), &jsonResult); err != nil {
		return nil, err
	}

	for database, collections := range jsonResult {
		if collections == nil {
			log.Info("Collecting all collections for database '%s'", database)
			filterMap.Filters[database] = nil
			continue
		}
		filterMap.Filters[database] = make(map[string]struct{})
		for _, collection := range collections {
			filterMap.Filters[database][collection] = struct{}{}
		}
	}

	return &filterMap, nil
}

// CheckFilter takes a database name and collection name and checks them against the whitelist
func (dbFilter *DatabaseFilter) CheckFilter(dbName, collectionName string) bool {
	// no filter, no whitelist
	if dbFilter.CollectAll {
		return true
	}

	collections, ok := dbFilter.Filters[dbName]
	if !ok {
		// database not contained in filter, don't collect
		return false
	}

	if collections == nil || collectionName == "" {
		// either collecting all collections in this database, or we're only filtering at the database level for this call.
		return true
	}

	// return whether or not collection is included in filter
	_, ok = collections[collectionName]
	return ok
}

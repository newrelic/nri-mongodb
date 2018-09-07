package entities

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// hostCollector is a base collector for any entity that represents a specific host
type hostCollector struct {
	defaultCollector
}

type cmdLineOpts struct {
	Argv   []string
	Parsed map[string]interface{}
	Ok     float64
}

// collectInventory collects all the inventory for a given host.
// This needs to be called from the parent object (Mongod/Mongos/Config)
// so that the integration.Entity is from the parent.
func (hc *hostCollector) collectInventory(e *integration.Entity) {
	hc.populateParameters(e)
	hc.populateCmdLineOpts(e)
}

func (hc *hostCollector) populateCmdLineOpts(entity *integration.Entity) {
	var cmdOpts cmdLineOpts

	if err := hc.session.DB("admin").Run(cmd{"getCmdLineOpts": 1}, &cmdOpts); err != nil {
		log.Error("Error calling getCmdLineOpts for [%s]: %v", hc.name, err)
		return
	}
	if cmdOpts.Ok == 1 {
		for i, arg := range cmdOpts.Argv {
			setInventoryItem(entity, "commandline", "argv", strconv.Itoa(i), arg)
		}
		addInventoryMap(entity, "configuration", cmdOpts.Parsed, false)
	}
}

func (hc *hostCollector) populateParameters(entity *integration.Entity) {
	var params map[string]interface{}

	if err := hc.session.DB("admin").Run(cmd{"getParameter": "*"}, &params); err != nil {
		log.Error("Error calling getParameter for [%s]: %v", hc.name, err)
		return
	}
	ok, exists := params["ok"]
	if exists && ok.(float64) == 1 {
		addInventoryMap(entity, "parameter", params, true)
	}
}

func addInventoryItem(entity *integration.Entity, category, key string, value interface{}, keyPrefix ...string) {
	switch v := value.(type) {
	case []interface{}:
		addInventoryArray(entity, category, key, v, keyPrefix...)
	case map[string]interface{}:
		addInventoryMap(entity, category, v, false, append(keyPrefix, key)...)
	case string:
		if v != "" {
			setInventoryItem(entity, category, key, "value", v, keyPrefix...)
		}
	default:
		setInventoryItem(entity, category, key, "value", v, keyPrefix...)
	}
}

func addInventoryArray(entity *integration.Entity, category, key string, value []interface{}, keyPrefix ...string) {
	if len(value) > 0 {
		values := make([]string, len(value))
		for i, v := range value {
			values[i] = fmt.Sprint(v)
		}
		joined := strings.Join(values, ", ")
		setInventoryItem(entity, category, key, "value", joined, keyPrefix...)
	}
}

func addInventoryMap(entity *integration.Entity, category string, value map[string]interface{}, isRoot bool, keyPrefix ...string) {
	for k, v := range value {
		if isInventoryKey(k, isRoot) {
			addInventoryItem(entity, category, k, v, keyPrefix...)
		}
	}
}

var ignoredInventoryKeys = map[string]bool{
	"operationTime": true,
	"ok":            true,
}

func isInventoryKey(key string, isRoot bool) bool {
	if !isRoot {
		return true
	}
	if strings.HasPrefix(key, "$") {
		return false
	}
	if ignore, ok := ignoredInventoryKeys[key]; ok && ignore {
		return false
	}
	return true
}

func setInventoryItem(entity *integration.Entity, category, key, field string, value interface{}, keyPrefix ...string) {
	key = strings.Join(append(keyPrefix, key), ".")
	if category != "" {
		key = category + "/" + key
	}
	if err := entity.SetInventoryItem(key, field, value); err != nil {
		log.Warn("Error setting inventory item [%s] %s=%v, %v", key, field, value, err)
	}
}

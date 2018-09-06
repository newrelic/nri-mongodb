package metrics

// Top is a storage struct for the top command
type Top struct {
	Totals map[string]TopRecords `bson:"totals"`
}

// TopRecords is a storage struct
type TopRecords struct {
	Total     TopRecordsTotal
	ReadLock  TopRecordsReadLock
	WriteLock TopRecordsWriteLock
	Queries   TopRecordsQueries
	Getmore   TopRecordsGetmore
	Insert    TopRecordsInsert
	Update    TopRecordsUpdate
	Remove    TopRecordsRemove
	Commands  TopRecordsCommands
}

// TopRecordsTotal is a storage struct
type TopRecordsTotal struct {
	Time  *int `bson:"time" metric_name:"usage.totalInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.totalPerSecond" source_type:"rate"`
}

// TopRecordsReadLock is a storage struct
type TopRecordsReadLock struct {
	Time  *int `bson:"time" metric_name:"usage.readLockInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.readLockPerSecond" source_type:"rate"`
}

// TopRecordsWriteLock is a storage struct
type TopRecordsWriteLock struct {
	Time  *int `bson:"time" metric_name:"usage.writeLockInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.writeLockPerSecond" source_type:"rate"`
}

// TopRecordsQueries is a storage struct
type TopRecordsQueries struct {
	Time  *int `bson:"time" metric_name:"usage.queriesInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.queriesPerSecond" source_type:"rate"`
}

// TopRecordsGetmore is a storage struct
type TopRecordsGetmore struct {
	Time  *int `bson:"time" metric_name:"usage.getmoreInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.getmorePerSecond" source_type:"rate"`
}

// TopRecordsInsert is a storage struct
type TopRecordsInsert struct {
	Time  *int `bson:"time" metric_name:"usage.insertInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.insertPerSecond" source_type:"rate"`
}

// TopRecordsUpdate is a storage struct
type TopRecordsUpdate struct {
	Time  *int `bson:"time" metric_name:"usage.updateInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.updatePerSecond" source_type:"rate"`
}

// TopRecordsRemove is a storage struct
type TopRecordsRemove struct {
	Time  *int `bson:"time" metric_name:"usage.removeInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.removePerSecond" source_type:"rate"`
}

// TopRecordsCommands is a storage struct
type TopRecordsCommands struct {
	Time  *int `bson:"time" metric_name:"usage.commandsInMilliseconds" source_type:"gauge"`
	Count *int `bson:"count" metric_name:"usage.commandsPerSecond" source_type:"rate"`
}

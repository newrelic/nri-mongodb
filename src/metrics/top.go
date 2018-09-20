package metrics

// Top is a storage struct for the top command
type Top struct {
	Totals map[string]TopRecords `bson:"totals"`
}

// TopRecords is a storage struct
type TopRecords struct {
	Total     *TopRecordsTotal     `bson:"total"`
	ReadLock  *TopRecordsReadLock  `bson:"readLock"`
	WriteLock *TopRecordsWriteLock `bson:"writeLock"`
	Queries   *TopRecordsQueries   `bson:"queries"`
	Getmore   *TopRecordsGetmore   `bson:"getmore"`
	Insert    *TopRecordsInsert    `bson:"insert"`
	Update    *TopRecordsUpdate    `bson:"update"`
	Remove    *TopRecordsRemove    `bson:"remove"`
	Commands  *TopRecordsCommands  `bson:"commands"`
}

// TopRecordsTotal is a storage struct
type TopRecordsTotal struct {
	Time  *int `bson:"time" metric_name:"usage.totalInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.totalPerSecond" source_type:"rate"`
}

// TopRecordsReadLock is a storage struct
type TopRecordsReadLock struct {
	Time  *int `bson:"time" metric_name:"usage.readLockInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.readLockPerSecond" source_type:"rate"`
}

// TopRecordsWriteLock is a storage struct
type TopRecordsWriteLock struct {
	Time  *int `bson:"time" metric_name:"usage.writeLockInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.writeLockPerSecond" source_type:"rate"`
}

// TopRecordsQueries is a storage struct
type TopRecordsQueries struct {
	Time  *int `bson:"time" metric_name:"usage.queriesInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.queriesPerSecond" source_type:"rate"`
}

// TopRecordsGetmore is a storage struct
type TopRecordsGetmore struct {
	Time  *int `bson:"time" metric_name:"usage.getmoreInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.getmorePerSecond" source_type:"rate"`
}

// TopRecordsInsert is a storage struct
type TopRecordsInsert struct {
	Time  *int `bson:"time" metric_name:"usage.insertInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.insertPerSecond" source_type:"rate"`
}

// TopRecordsUpdate is a storage struct
type TopRecordsUpdate struct {
	Time  *int `bson:"time" metric_name:"usage.updateInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.updatePerSecond" source_type:"rate"`
}

// TopRecordsRemove is a storage struct
type TopRecordsRemove struct {
	Time  *int `bson:"time" metric_name:"usage.removeInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.removePerSecond" source_type:"rate"`
}

// TopRecordsCommands is a storage struct
type TopRecordsCommands struct {
	Time  *int `bson:"time" metric_name:"usage.commandsInMillisecondsPerSecond" source_type:"rate"`
	Count *int `bson:"count" metric_name:"usage.commandsPerSecond" source_type:"rate"`
}

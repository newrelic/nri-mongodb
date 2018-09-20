package metrics

// ServerStatus is a storage struct
type ServerStatus struct {
	Host                      *string                         `bson:"host"`
	Version                   *string                         `bson:"version"`
	Process                   *string                         `bson:"process"`
	PID                       *int                            `bson:"pid"`
	Uptime                    *int                            `bson:"uptime"`
	Asserts                   *ServerStatusAsserts            `bson:"asserts"`
	Connections               *ServerStatusConnections        `bson:"connections"`
	LogicalSessionRecordCache *ServerStatusLSRC               `bson:"logicalSessionRecordCache"`
	Network                   *ServerStatusNetwork            `bson:"network"`
	Opcounters                *ServerStatusOpcounters         `bson:"opcounters"`
	OpcountersRepl            *ServerStatusOpcountersRepl     `bson:"opcountersRepl"`
	Mem                       *ServerStatusMem                `bson:"mem"`
	Metrics                   *ServerStatusMetrics            `bson:"metrics"`
	BackgroundFlushing        *ServerStatusBackgroundFlushing `bson:"backgroundFlushing"`
	GlobalLock                *ServerStatusGlobalLock         `bson:"globalLock"`
	ExtraInfo                 *ServerStatusExtraInfo          `bson:"extra_info"`
	WiredTiger                *ServerStatusWiredTiger         `bson:"wiredTiger"`
	Locks                     *ServerStatusLocks              `bson:"locks"`
	Dur                       *ServerStatusDur                `bson:"dur"`
}

// ServerStatusAsserts is a storage struct
type ServerStatusAsserts struct {
	Regular   *int `bson:"regular"   metric_name:"asserts.regularPerSecond"   source_type:"rate"`
	Warning   *int `bson:"warning"   metric_name:"asserts.warningPerSecond"   source_type:"rate"`
	Msg       *int `bson:"msg"       metric_name:"asserts.messagesPerSecond"  source_type:"rate"`
	User      *int `bson:"user"      metric_name:"asserts.userPerSecond"      source_type:"rate"`
	Rollovers *int `bson:"rollovers" metric_name:"asserts.rolloversPerSecond" source_type:"rate"`
}

// ServerStatusConnections is a storage struct
type ServerStatusConnections struct {
	Current      *int `bson:"current"      metric_name:"connections.current"      source_type:"gauge"`
	Available    *int `bson:"available"    metric_name:"connections.available"    source_type:"gauge"`
	TotalCreated *int `bson:"totalCreated" metric_name:"connections.totalCreated" source_type:"gauge"`
}

// ServerStatusLSRC is a storage struct
type ServerStatusLSRC struct {
	ActiveSessionsCount                       *int `bson:"activeSessionsCount"`
	SessionsCollectionJobCount                *int `bson:"sessionsCollectionJobCount"`
	LastSessionsCollectionJobDurationMillis   *int `bson:"lastSessionsCollectionJobDurationMillis"`
	LastSessionsCollectionJobEntriesRefreshed *int `bson:"lastSessionsCollectionJobEntriesRefreshed"`
	LastSessionsCollectionJobEntriesClosed    *int `bson:"lastSessionCollectionJobEntriesClosed"`
	TransactionReaperJobCount                 *int `bson:"transactionReaperJobCount"`
	LastTransactionReaperJobDurationMillis    *int `bson:"lastTransactionReaperJobDurationMillis"`
	LastTransactionReaperJobEntriesCleanedUp  *int `bson:"lastTransactionReaperJobEntriesCleanedUp"`
}

// ServerStatusNetwork is a storage struct
type ServerStatusNetwork struct {
	BytesIn     *int `bson:"bytesIn"     metric_name:"network.bytesInPerSecond"  source_type:"rate"`
	BytesOut    *int `bson:"bytesOut"    metric_name:"network.bytesOutPerSecond" source_type:"rate"`
	NumRequests *int `bson:"numRequests" metric_name:"network.requestsPerSecond" source_type:"rate"`
}

// ServerStatusOpcounters is a storage struct
type ServerStatusOpcounters struct {
	Insert  *int `bson:"insert"  metric_name:"opcounters.insertPerSecond"  source_type:"rate"`
	Query   *int `bson:"query"   metric_name:"opcounters.queryPerSecond"   source_type:"rate"`
	Update  *int `bson:"update"  metric_name:"opcounters.updatePerSecond"  source_type:"rate"`
	Delete  *int `bson:"delete"  metric_name:"opcounters.deletePerSecond"  source_type:"rate"`
	Getmore *int `bson:"getmore" metric_name:"opcounters.getmorePerSecond" source_type:"rate"`
	Command *int `bson:"command" metric_name:"opcounters.commandPerSecond" source_type:"rate"`
}

// ServerStatusOpcountersRepl is a storage struct
type ServerStatusOpcountersRepl struct {
	Insert  *int `bson:"insert"  metric_name:"opcountersrepl.insertPerSecond"  source_type:"rate"`
	Query   *int `bson:"query"   metric_name:"opcountersrepl.queryPerSecond"   source_type:"rate"`
	Update  *int `bson:"update"  metric_name:"opcountersrepl.updatePerSecond"  source_type:"rate"`
	Delete  *int `bson:"delete"  metric_name:"opcountersrepl.deletePerSecond"  source_type:"rate"`
	Getmore *int `bson:"getmore" metric_name:"opcountersrepl.getmorePerSecond" source_type:"rate"`
	Command *int `bson:"command" metric_name:"opcountersrepl.commandPerSecond" source_type:"rate"`
}

// ServerStatusMem is a storage struct
type ServerStatusMem struct {
	Bits              *int `bson:"bits"              metric_name:"mem.bits"                     source_type:"gauge"`
	Resident          *int `bson:"resident"          metric_name:"mem.residentInBytes"          source_type:"gauge"`
	Virtual           *int `bson:"virtual"           metric_name:"mem.virtualInBytes"           source_type:"gauge"`
	Mapped            *int `bson:"mapped"            metric_name:"mem.mappedInBytes"            source_type:"gauge"`
	MappedWithJournal *int `bson:"mappedWithJournal" metric_name:"mem.mappedWithJournalInBytes" source_type:"gauge"`
}

// ServerStatusMetrics is a storage struct
type ServerStatusMetrics struct {
	Cursor        *ServerStatusMetricsCursor        `bson:"cursor"`
	Commands      *ServerStatusMetricsCommands      `bson:"commands"`
	Document      *ServerStatusMetricsDocument      `bson:"document"`
	GetLastError  *ServerStatusMetricsGetLastError  `bson:"getLastError"`
	Operation     *ServerStatusMetricsOperation     `bson:"operation"`
	QueryExecutor *ServerStatusMetricsQueryExecutor `bson:"queryExecutor"`
	Record        *ServerStatusMetricsRecord        `bson:"record"`
	Repl          *ServerStatusMetricsRepl          `bson:"repl"`
	TTL           *ServerStatusMetricsTTL           `bson:"ttl"`
}

// ServerStatusMetricsCursor is a storage struct
type ServerStatusMetricsCursor struct {
	TimedOutPerSecond *int                           `bson:"timedOut" metric_name:"cursor.timedOutPerSecond" source_type:"rate"`
	Open              *ServerStatusMetricsCursorOpen `bson:"open"`
}

// ServerStatusMetricsCursorOpen is a storage struct
type ServerStatusMetricsCursorOpen struct {
	Total     *int `bson:"total"  metric_name:"cursor.openTotal"  source_type:"gauge"`
	Pinned    *int `bson:"pinned" metric_name:"cursor.openPinned" source_type:"gauge"`
	NoTimeout *int `bson:"noTimeout" metric_name:"cursor.openNoTimeout" source_type:"gauge"`
}

// ServerStatusMetricsCommands is a storage struct
type ServerStatusMetricsCommands struct {
	Count         *ServerStatusMetricsCommandCount         `bson:"count"`
	CreateIndexes *ServerStatusMetricsCommandCreateIndexes `bson:"createIndexes"`
	Delete        *ServerStatusMetricsCommandDelete        `bson:"delete"`
	Eval          *ServerStatusMetricsCommandEval          `bson:"eval"`
	FindAndModify *ServerStatusMetricsCommandFindAndModify `bson:"findAndModify"`
	Insert        *ServerStatusMetricsCommandInsert        `bson:"insert"`
	Update        *ServerStatusMetricsCommandUpdate        `bson:"update"`
}

// ServerStatusMetricsDocument is a storage struct
type ServerStatusMetricsDocument struct {
	Deleted  *int `bson:"deleted"  metric_name:"document.deletedPerSecond"  source_type:"rate"`
	Inserted *int `bson:"inserted" metric_name:"document.insertedPerSecond" source_type:"rate"`
	Returned *int `bson:"returned" metric_name:"document.returnedPerSecond" source_type:"rate"`
	Updated  *int `bson:"updated"  metric_name:"document.updatedPerSecond"  source_type:"rate"`
}

// ServerStatusMetricsCommandCount is a storage struct
type ServerStatusMetricsCommandCount struct {
	Failed *int `bson:"failed" metric_name:"commands.countFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.countPerSecond"       source_type:"rate"`
}

// ServerStatusMetricsCommandCreateIndexes is a storage struct
type ServerStatusMetricsCommandCreateIndexes struct {
	Failed *int `bson:"failed" metric_name:"commands.createIndexesFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.createIndexesPerSecond"       source_type:"rate"`
}

// ServerStatusMetricsCommandDelete is a storage struct
type ServerStatusMetricsCommandDelete struct {
	Failed *int `bson:"failed" metric_name:"commands.deleteFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.deletePerSecond"       source_type:"rate"`
}

// ServerStatusMetricsCommandEval is a storage struct
type ServerStatusMetricsCommandEval struct {
	Failed *int `bson:"failed" metric_name:"commands.evalFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.evalPerSecond"       source_type:"rate"`
}

// ServerStatusMetricsCommandFindAndModify is a storage struct
type ServerStatusMetricsCommandFindAndModify struct {
	Failed *int `bson:"failed" metric_name:"commands.findAndModifyFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.findAndModifyPerSecond"           source_type:"rate"`
}

// ServerStatusMetricsCommandInsert is a storage struct
type ServerStatusMetricsCommandInsert struct {
	Failed *int `bson:"failed" metric_name:"commands.insertFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.insertPerSecond"       source_type:"rate"`
}

// ServerStatusMetricsCommandUpdate is a storage struct
type ServerStatusMetricsCommandUpdate struct {
	Failed *int `bson:"failed" metric_name:"commands.updateFailedPerSecond" source_type:"rate"`
	Total  *int `bson:"total"  metric_name:"commands.updatePerSecond"       source_type:"rate"`
}

// ServerStatusMetricsGetLastError is a storage struct
type ServerStatusMetricsGetLastError struct {
	Wtime     *ServerStatusMetricsGetLastErrorWtime `bson:"wtime"`
	Wtimeouts *int                                  `bson:"wtimeouts" metric_name:"getlasterror.wtimeoutsPerSecond" source_type:"rate"`
}

// ServerStatusMetricsGetLastErrorWtime is a storage struct
type ServerStatusMetricsGetLastErrorWtime struct {
	TotalMillis *int `bson:"totalMillis" metric_name:"getlasterror.wtimeMillisPerSecond" source_type:"rate"`
}

// ServerStatusMetricsOperation is a storage struct
type ServerStatusMetricsOperation struct {
	ScanAndOrder   *int `bson:"scanAndOrder"   metric_name:"operation.scanAndOrderPerSecond"   source_type:"rate"`
	WriteConflicts *int `bson:"writeConflicts" metric_name:"operation.writeConflictsPerSecond" source_type:"rate"`
	Idhack         *int `bson:"idhack"         metric_name:"operation.idhackPerSecond"         source_type:"rate"`
	Fastmod        *int `bson:"fastmod"        metric_name:"operation.fastmodPerSecond"        source_type:"rate"`
}

// ServerStatusMetricsQueryExecutor is a storage struct
type ServerStatusMetricsQueryExecutor struct {
	Scanned *int `bson:"scanned" metric_name:"queryexecutor.scannedPerSecond" source_type:"rate"`
}

// ServerStatusMetricsRecord is a storage struct
type ServerStatusMetricsRecord struct {
	Moves *int `bson:"moves" metric_name:"record.movesPerSecond" source_type:"rate"`
}

// ServerStatusMetricsRepl is a storage struct
type ServerStatusMetricsRepl struct {
	Apply   *ServerStatusMetricsReplApply   `bson:"apply"`
	Buffer  *ServerStatusMetricsReplBuffer  `bson:"buffer"`
	Network *ServerStatusMetricsReplNetwork `bson:"network"`
	Preload *ServerStatusMetricsReplPreload `bson:"preload"`
	Oplog   *ServerStatusMetricsReplOplog   `bson:"oplog"`
}

// ServerStatusMetricsReplApply is a storage struct
type ServerStatusMetricsReplApply struct {
	Ops     *int                                 `bson:"ops" metric_name:"repl.apply.operationsPerSecond" source_type:"rate"`
	Batches *ServerStatusMetricsReplApplyBatches `bson:"batches"`
}

// ServerStatusMetricsReplApplyBatches is a storage struct
type ServerStatusMetricsReplApplyBatches struct {
	Num *int `bson:"num" metric_name:"repl.apply.batchesPerSecond" source_type:"rate"`
}

// ServerStatusMetricsReplBuffer is a storage struct
type ServerStatusMetricsReplBuffer struct {
	Count        *int `bson:"count"        metric_name:"repl.buffer.count"          source_type:"gauge"`
	MaxSizeBytes *int `bson:"maxSizeBytes" metric_name:"repl.buffer.maxSizeInBytes" source_type:"gauge"`
	SizeBytes    *int `bson:"sizeBytes"    metric_name:"repl.buffer.sizeInBytes"    source_type:"gauge"`
}

// ServerStatusMetricsReplNetwork is a storage struct
type ServerStatusMetricsReplNetwork struct {
	Bytes          *int                                    `bson:"bytes"          metric_name:"repl.network.bytesPerSecond"          source_type:"rate"`
	Ops            *int                                    `bson:"ops"            metric_name:"repl.network.operationPerSecond"      source_type:"rate"`
	ReadersCreated *int                                    `bson:"readersCreated" metric_name:"repl.network.readersCreatedPerSecond" source_type:"rate"`
	Getmores       *ServerStatusMetricsReplNetworkGetmores `bson:"getmores"`
}

// ServerStatusMetricsReplNetworkGetmores is a storage struct
type ServerStatusMetricsReplNetworkGetmores struct {
	Num *int `bson:"num" metric_name:"repl.network.getmoresPerSecond" source_type:"rate"`
}

// ServerStatusMetricsReplPreload is a storage struct
type ServerStatusMetricsReplPreload struct {
	Docs    *ServerStatusMetricsReplPreloadDocs    `bson:"docs"`
	Indexes *ServerStatusMetricsReplPreloadIndexes `bson:"indexes"`
}

// ServerStatusMetricsReplOplog is a storage struct
type ServerStatusMetricsReplOplog struct {
	Docs    *ServerStatusMetricsReplPreloadDocs    `bson:"docs"`
	Indexes *ServerStatusMetricsReplPreloadIndexes `bson:"indexes"`
}

// ServerStatusMetricsReplPreloadDocs is a storage struct
type ServerStatusMetricsReplPreloadDocs struct {
	Num         *int `bson:"num"         metric_name:"repl.docsLoadedPrefetch"        source_type:"gauge"`
	TotalMillis *int `bson:"totalMillis" metric_name:"repl.docsPreloadInMilliseconds" source_type:"gauge"`
}

// ServerStatusMetricsReplPreloadIndexes is a storage struct
type ServerStatusMetricsReplPreloadIndexes struct {
	Num         *int `bson:"num"         metric_name:"repl.indexLoadedPrefetch"        source_type:"gauge"`
	TotalMillis *int `bson:"totalMillis" metric_name:"repl.indexPreloadInMilliseconds" source_type:"gauge"`
}

// ServerStatusMetricsTTL is a storage struct
type ServerStatusMetricsTTL struct {
	DeletedDocuments *int `bson:"deletedDocuments" metric_name:"ttl.deleteDocumentsPerSecond" source_type:"rate"`
	Passes           *int `bson:"passes"           metric_name:"ttl.removeDocumentPerSecond"  source_type:"rate"`
}

// ServerStatusBackgroundFlushing is a storage struct
type ServerStatusBackgroundFlushing struct {
	Flushes   *int     `bson:"flushes"    metric_name:"flush.flushesDisk"           source_type:"gauge"`
	TotalMs   *float64 `bson:"total_ms"   metric_name:"flush.totalInMilliseconds"     source_type:"gauge"`
	AverageMs *float64 `bson:"average_ms" metric_name:"flush.averageInMilliseconds" source_type:"gauge"`
	LastMs    *float64 `bson:"last_ms"    metric_name:"flush.lastInMilliseconds"    source_type:"gauge"`
}

// ServerStatusGlobalLock is a storage struct
type ServerStatusGlobalLock struct {
	TotalTime     *float32                             `bson:"totalTime" metric_name:"globallock.totaltime" source_type:"gauge"`
	CurrentQueue  *ServerStatusGlobalLockCurrentQueue  `bson:"currentQueue"`
	ActiveClients *ServerStatusGlobalLockActiveClients `bson:"activeClients"`
}

// ServerStatusGlobalLockCurrentQueue is a storage struct
type ServerStatusGlobalLockCurrentQueue struct {
	Total   *int `bson:"total" metric_name:"globallock.currentQueueTotal"   source_type:"gauge"`
	Readers *int `bson:"readers"   metric_name:"globallock.currentQueueReaders" source_type:"gauge"`
	Writers *int `bson:"writers"   metric_name:"globallock.currentQueueWriters" source_type:"gauge"`
}

// ServerStatusGlobalLockActiveClients is a storage struct
type ServerStatusGlobalLockActiveClients struct {
	Total   *int `bson:"total" metric_name:"globallock.activeClientsTotal"   source_type:"gauge"`
	Readers *int `bson:"readers"   metric_name:"globallock.activeClientsReaders" source_type:"gauge"`
	Writers *int `bson:"writers"   metric_name:"globallock.activeClientsWriters" source_type:"gauge"`
}

// ServerStatusExtraInfo is a storage struct
type ServerStatusExtraInfo struct {
	PageFaults *int `bson:"page_faults" metric_name:"pageFaultsPerSecond" source_type:"rate"`
}

// ServerStatusWiredTiger is a storage struct
type ServerStatusWiredTiger struct {
	Cache                  *ServerStatusWiredTigerCache                  `bson:"cache"`
	ConcurrentTransactions *ServerStatusWiredTigerConcurrentTransactions `bson:"concurrentTransactions"`
}

// ServerStatusWiredTigerCache is a storage struct
type ServerStatusWiredTigerCache struct {
	Size                   *int `bson:"bytes currently in the cache"                                 metric_name:"wiredtiger.cacheInBytes"                                 source_type:"gauge"`
	FailedEvictions        *int `bson:"failed eviction of pages that exceeded the in-memory maximum count" metric_name:"wiredtiger.failedEvictionsPagesPerSecond"                source_type:"rate"`
	PageSplits             *int `bson:"in-memory page splits"                                        metric_name:"wiredtiger.cacheInMemoryPageSplits"                      source_type:"gauge"`
	MaxSize                *int `bson:"maximum bytes configured"                                     metric_name:"wiredtiger.cacheMaxInBytes"                              source_type:"gauge"`
	MaxPageSize            *int `bson:"maximum page size at eviction"                                metric_name:"wiredtiger.cacheMaxPageSizeEvictionInBytes"              source_type:"gauge"`
	ModifiedPagesEvicted   *int `bson:"modified pages evicted"                                       metric_name:"wiredtiger.cacheModifiedPagesEvicted"                    source_type:"gauge"`
	PagesHeld              *int `bson:"pages currently held in the cache"                            metric_name:"wiredtiger.cachePagesHeld"                               source_type:"gauge"`
	PagesEvictedThreads    *int `bson:"pages evicted by application threads"                         metric_name:"wiredtiger.cachePagesEvictedApplicationThreadsPerSecond" source_type:"gauge"`
	PagesEvictedMax        *int `bson:"pages evicted because they exceeded the in-memory maximum count"    metric_name:"wiredtiger.cachePagesEvictedInMemoryMaxPerSecond"        source_type:"rate"`
	DirtyData              *int `bson:"tracked dirty bytes in the cache"                             metric_name:"wiredtiger.cacheDirtyDataInBytes"                        source_type:"gauge"`
	UnmodifiedPagesEvicted *int `bson:"unmodified pages evicted"                                     metric_name:"wiredtiger.cacheUnmodifiedPagesEvicted"                  source_type:"gauge"`
}

// ServerStatusWiredTigerConcurrentTransactions is a storage struct
type ServerStatusWiredTigerConcurrentTransactions struct {
	Write *ServerStatusWiredTigerConcurrentTransactionsWrite `bson:"write"`
	Read  *ServerStatusWiredTigerConcurrentTransactionsRead  `bson:"read"`
}

// ServerStatusWiredTigerConcurrentTransactionsWrite is a storage struct
type ServerStatusWiredTigerConcurrentTransactionsWrite struct {
	Out          *int `bson:"out"          metric_name:"wiredtiger.concurrentTransactions.WriteRemaining" source_type:"gauge"`
	Available    *int `bson:"available"    metric_name:"wiredtiger.concurrentTransactions.WriteAvailable" source_type:"gauge"`
	TotalTickets *int `bson:"totalTickets" metric_name:"wiredtiger.concurrentTransactions.WriteTotal"     source_type:"gauge"`
}

// ServerStatusWiredTigerConcurrentTransactionsRead is a storage struct
type ServerStatusWiredTigerConcurrentTransactionsRead struct {
	Out          *int `bson:"out"          metric_name:"wiredtiger.concurrentTransactions.ReadRemaining" source_type:"gauge"`
	Available    *int `bson:"available"    metric_name:"wiredtiger.concurrentTransactions.ReadAvailable" source_type:"gauge"`
	TotalTickets *int `bson:"totalTickets" metric_name:"wiredtiger.concurrentTransactions.ReadTotal"     source_type:"gauge"`
}

// ServerStatusLocks is a storage struct
type ServerStatusLocks struct {
	Collection    *ServerStatusLocksCollection    `bson:"Collection"`
	Database      *ServerStatusLocksDatabase      `bson:"Database"`
	Global        *ServerStatusLocksGlobal        `bson:"Global"`
	Metadata      *ServerStatusLocksMetadata      `bson:"Metadata"`
	MMAPV1Journal *ServerStatusLocksMMAPV1Journal `bson:"MMAPV1Journal"`
	Oplog         *ServerStatusLocksOplog         `bson:"oplog"`
}

// ServerStatusLocksCollection is a storage struct
type ServerStatusLocksCollection struct {
	AcquireCount        *ServerStatusLocksCollectionAcquireCount        `bson:"acquireCount"`
	AcquireWaitCount    *ServerStatusLocksCollectionAcquireWaitCount    `bson:"acquireWaitCount"`
	TimeAcquiringMicros *ServerStatusLocksCollectionTimeAcquiringMicros `bson:"timeAcquiringMicros"`
}

// ServerStatusLocksCollectionAcquireCount is a storage struct
type ServerStatusLocksCollectionAcquireCount struct {
	Shared          *int `bson:"R" metric_name:"locks.collectionAcquireShared" source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.collectionAcquireExclusive" source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.collectionAcquireIntentShared" source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.collectionAcquireIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksCollectionAcquireWaitCount is a storage struct
type ServerStatusLocksCollectionAcquireWaitCount struct {
	Shared    *int `bson:"R" metric_name:"locks.collectionAcquireWaitCountShared"    source_type:"gauge"`
	Exclusive *int `bson:"W" metric_name:"locks.collectionAcquireWaitCountExclusive" source_type:"gauge"`
}

// ServerStatusLocksCollectionTimeAcquiringMicros is a storage struct
type ServerStatusLocksCollectionTimeAcquiringMicros struct {
	Shared    *int `bson:"R" metric_name:"locks.collectionTimeAcquiringMicrosShared"    source_type:"gauge"`
	Exclusive *int `bson:"W" metric_name:"locks.collectionTimeAcquiringMicrosExclusive" source_type:"gauge"`
}

// ServerStatusLocksDatabase is a storage struct
type ServerStatusLocksDatabase struct {
	AcquireCount        *ServerStatusLocksDatabaseAcquireCount        `bson:"acquireCount"`
	AcquireWaitCount    *ServerStatusLocksDatabaseAcquireWaitCount    `bson:"acquireWaitCount"`
	TimeAcquiringMicros *ServerStatusLocksDatabaseTimeAcquiringMicros `bson:"timeAcquiringMicros"`
}

// ServerStatusLocksDatabaseAcquireCount is a storage struct
type ServerStatusLocksDatabaseAcquireCount struct {
	Shared          *int `bson:"R" metric_name:"locks.databaseAcquireShared" source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.databaseAcquireExclusive" source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.databaseAcquireIntentShared" source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.databaseAcquireIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksDatabaseAcquireWaitCount is a storage struct
type ServerStatusLocksDatabaseAcquireWaitCount struct {
	Shared          *int `bson:"R" metric_name:"locks.databaseAcquireWaitShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.databaseAcquireWaitExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.databaseAcquireWaitIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.databaseAcquireWaitIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksDatabaseTimeAcquiringMicros is a storage struct
type ServerStatusLocksDatabaseTimeAcquiringMicros struct {
	Shared          *int `bson:"R" metric_name:"locks.databaseTimeAcquiringMicrosShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.databaseTimeAcquiringMicrosExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.databaseTimeAcquiringMicrosIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.databaseTimeAcquiringMicrosIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksGlobal is a storage struct
type ServerStatusLocksGlobal struct {
	AcquireCount        *ServerStatusLocksGlobalAcquireCount        `bson:"acquireCount"`
	AcquireWaitCount    *ServerStatusLocksGlobalAcquireWaitCount    `bson:"acquireWaitCount"`
	TimeAcquiringMicros *ServerStatusLocksGlobalTimeAcquiringMicros `bson:"timeAcquiringMicros"`
}

// ServerStatusLocksGlobalAcquireCount is a storage struct
type ServerStatusLocksGlobalAcquireCount struct {
	Shared          *int `bson:"R" metric_name:"locks.globalAcquireShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.globalAcquireExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.globalAcquireIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.globalAcquireIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksGlobalAcquireWaitCount is a storage struct
type ServerStatusLocksGlobalAcquireWaitCount struct {
	Shared          *int `bson:"R" metric_name:"locks.globalAcquireWaitShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.globalAcquireWaitExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.globalAcquireWaitIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.globalAcquireWaitIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksGlobalTimeAcquiringMicros is a storage struct
type ServerStatusLocksGlobalTimeAcquiringMicros struct {
	Shared          *int `bson:"R" metric_name:"locks.globalTimeAcquiringMicrosShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.globalTimeAcquiringMicrosExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.globalTimeAcquiringMicrosIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.globalTimeAcquiringMicrosIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksMetadata is a storage struct
type ServerStatusLocksMetadata struct {
	AcquireCount *ServerStatusLocksMetadataAcquireCount `bson:"acquireCount"`
}

// ServerStatusLocksMetadataAcquireCount is a storage struct
type ServerStatusLocksMetadataAcquireCount struct {
	Shared    *int `bson:"R" metric_name:"locks.metadataAcquireShared" source_type:"gauge"`
	Exclusive *int `bson:"W" metric_name:"locks.metadataAcquireExclusive" source_type:"gauge"`
}

// ServerStatusLocksMMAPV1Journal is a storage struct
type ServerStatusLocksMMAPV1Journal struct {
	AcquireCount        *ServerStatusLocksMMAPV1JournalAcquireCount        `bson:"acquireCount"`
	TimeAcquiringMicros *ServerStatusLocksMMAPV1JournalTimeAcquiringMicros `bson:"timeAcquiringMicros"`
}

// ServerStatusLocksMMAPV1JournalAcquireCount is a storage struct
type ServerStatusLocksMMAPV1JournalAcquireCount struct {
	Shared          *int `bson:"R" metric_name:"locks.mmapv1journalAcquireShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.mmapv1journalAcquireExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.mmapv1journalAcquireIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.mmapv1journalAcquireIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksOplog is a storage struct
type ServerStatusLocksOplog struct {
	AcquireCount        *ServerStatusLocksOplogAcquireCount        `bson:"acquireCount"`
	AcquireWaitCount    *ServerStatusLocksOplogAcquireCount        `bson:"acquireWaitCount"`
	TimeAcquiringMicros *ServerStatusLocksOplogTimeAcquiringMicros `bson:"timeAcquiringMicros"`
}

// ServerStatusLocksOplogAcquireCount is a storage struct
type ServerStatusLocksOplogAcquireCount struct {
	Shared          *int `bson:"R" metric_name:"locks.oplogAcquireShared"          source_type:"gauge"`
	Exclusive       *int `bson:"W" metric_name:"locks.oplogAcquireExclusive"       source_type:"gauge"`
	IntentShared    *int `bson:"r" metric_name:"locks.oplogAcquireIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.oplogAcquireIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksOplogAcquireWaitCount is a storage struct
type ServerStatusLocksOplogAcquireWaitCount struct {
	Shared          *int `bson:"R" metric_name:"locks.oplogAcquireWaitShared"          source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.oplogAcquireWaitIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksOplogTimeAcquiringMicros is a storage struct
type ServerStatusLocksOplogTimeAcquiringMicros struct {
	IntentShared    *int `bson:"r" metric_name:"locks.oplogTimeAcquiringMicrosIntentShared"    source_type:"gauge"`
	IntentExclusive *int `bson:"w" metric_name:"locks.oplogTimeAcquiringMicrosIntentExclusive" source_type:"gauge"`
}

// ServerStatusLocksMMAPV1JournalTimeAcquiringMicros is a storage struct
type ServerStatusLocksMMAPV1JournalTimeAcquiringMicros struct {
	Shared    *int `bson:"R" metric_name:"locks.mmapv1journalTimeAcquiringMicrosShared"          source_type:"gauge"`
	Exclusive *int `bson:"W" metric_name:"locks.mmapv1journalTimeAcquiringMicrosExclusive"       source_type:"gauge"`
}

// ServerStatusDur is a storage struct
type ServerStatusDur struct {
	Commits               *int                   `bson:"commits"                                       metric_name:"dur.commits"            source_type:"gauge"`
	CommitsInWriteLock    *int                   `bson:"commitsInWriteLock"                            metric_name:"dur.commitsInWriteLock" source_type:"gauge"`
	Compression           *float32               `bson:"compression"                                   metric_name:"dur.compression"        source_type:"gauge"`
	EarlyCommits          *int                   `bson:"earlyCommits"                                  metric_name:"dur.earlyCommits"       source_type:"gauge"`
	JournaledBytes        *float32               `metric_name:"dur.journaledInBytes"                   source_type:"gauge"`
	WriteToDataFilesBytes *float32               `metric_name:"dur.dataWrittenJournalDataFilesInBytes" source_type:"gauge"`
	JournaledMB           *float32               `bson:"journaledMB"`
	WriteToDataFilesMB    *float32               `bson:"writeToDataFilesMB"`
	TimeMS                *ServerStatusDurTimeMS `bson:"timeMs"`
}

// ServerStatusDurTimeMS is a storage struct
type ServerStatusDurTimeMS struct {
	Commits            *int `bson:"commits"            metric_name:"dur.commitsInMilliseconds"              source_type:"gauge"`
	CommitsInWriteLock *int `bson:"commitsInWriteLock" metric_name:"dur.commitsInWriteLockInMilliseconds"   source_type:"gauge"`
	Dt                 *int `bson:"dt"                 metric_name:"dur.timeCollectedCommitsInMilliseconds" source_type:"gauge"`
	PrepLogBuffer      *int `bson:"prepLogBuffer"      metric_name:"dur.preparingInMilliseconds"            source_type:"gauge"`
	RemapPrivateView   *int `bson:"remapPrivateView"   metric_name:"dur.remappingInMilliseconds"            source_type:"gauge"`
	WriteToDataFiles   *int `bson:"writeToDataFiles"   metric_name:"dur.writingDataFilesInMilliseconds"     source_type:"gauge"`
	WriteToJournal     *int `bson:"writeToJournal"     metric_name:"dur.writingJournalInMilliseconds"       source_type:"gauge"`
}

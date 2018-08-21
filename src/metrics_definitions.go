package main

type serverStatus struct {
	Host                      string
	Version                   string
	Process                   string
	PID                       int
	Uptime                    int
	Asserts                   serverStatusAsserts
	Connections               serverStatusConnections
	LogicalSessionRecordCache serverStatusLSRC
	Network                   serverStatusNetwork
	Opcounters                serverStatusOpcounters
	OpcountersRepl            serverStatusOpcountersRepl
	Mem                       serverStatusMem
	Metrics                   serverStatusMetrics
	BackgroundFlushing        serverStatusBackgroundFlushing
	GlobalLock                serverStatusGlobalLock
	ExtraInfo                 serverStatusExtraInfo
	WiredTiger                serverStatusWiredTiger
	Locks                     serverStatusLocks
	Dur                       serverStatusDur
}

type serverStatusAsserts struct {
	Regular   int `metric_name:"asserts.regularPerSecond"   source_type:"rate"`
	Warning   int `metric_name:"asserts.warningPerSecond"   source_type:"rate"`
	Msg       int `metric_name:"asserts.messagesPerSecond"  source_type:"rate"`
	User      int `metric_name:"asserts.userPerSecond"      source_type:"rate"`
	Rollovers int `metric_name:"asserts.rolloversPerSecond" source_type:"rate"`
}

type serverStatusConnections struct {
	Current      int `metric_name:"connections.current"      source_type:"gauge"`
	Available    int `metric_name:"connections.available"    source_type:"gauge"`
	TotalCreated int `metric_name:"connections.totalCreated" source_type:"gauge"`
}

type serverStatusLSRC struct {
	ActiveSessionsCount                       int
	SessionsCollectionJobCount                int
	LastSessionsCollectionJobDurationMillis   int
	LastSessionsCollectionJobEntriesRefreshed int
	LastSessionsCollectionJobEntriesClosed    int
	TransactionReapderJobCount                int
	LastTransactionReaperJobDurationMillis    int
	LastTransactionReaperJobEntriesCleanedUp  int
}

type serverStatusNetwork struct {
	BytesIn     int `metric_name:"network.bytesIn"  source_type:"gauge"`
	BytesOut    int `metric_name:"network.bytesOut" source_type:"gauge"`
	NumRequests int `metric_name:"network.requests" source_type:"gauge"`
}

type serverStatusOpcounters struct {
	Insert  int `metric_name:"opcounters.insertPerSecond"  source_type:"rate"`
	Query   int `metric_name:"opcounters.queryPerSecond"   source_type:"rate"`
	Update  int `metric_name:"opcounters.updatePerSecond"  source_type:"rate"`
	Delete  int `metric_name:"opcounters.deletePerSecond"  source_type:"rate"`
	Getmore int `metric_name:"opcounters.getmorePerSecond" source_type:"rate"`
	Command int `metric_name:"opcounters.commandPerSecond" source_type:"rate"`
}

type serverStatusOpcountersRepl struct {
	Insert  int `metric_name:"opcountersrepl.insertPerSecond"  source_type:"rate"`
	Query   int `metric_name:"opcountersrepl.queryPerSecond"   source_type:"rate"`
	Update  int `metric_name:"opcountersrepl.updatePerSecond"  source_type:"rate"`
	Delete  int `metric_name:"opcountersrepl.deletePerSecond"  source_type:"rate"`
	Getmore int `metric_name:"opcountersrepl.getmorePerSecond" source_type:"rate"`
	Command int `metric_name:"opcountersrepl.commandPerSecond" source_type:"rate"`
}

type serverStatusMem struct {
	Bits              int `metric_name:"mem.bits"                     source_type:"gauge"`
	Resident          int `metric_name:"mem.residentInBytes"          source_type:"gauge"`
	Virtual           int `metric_name:"mem.virtualInBytes"           source_type:"gauge"`
	Mapped            int `metric_name:"mem.mappedInBytes"            source_type:"gauge"`
	MappedWithJournal int `metric_name:"mem.mappedWithJournalInBytes" source_type:"gauge"`
}

type serverStatusMetrics struct {
	Cursor        serverStatusMetricsCursor
	Commands      serverStatusMetricsCommands
	GetLastError  serverStatusMetricsGetLastError
	Operation     serverStatusMetricsOperation
	QueryExecutor serverStatusMetricsQueryExecutor
	Record        serverStatusMetricsRecord
	Repl          serverStatusMetricsRepl
	Ttl           serverStatusMetricsTtl
}

type serverStatusMetricsCursor struct {
	TimedOut int `metric_name:"cursor.timedOutPerSecond" source_type:"rate"`
	Open     serverStatusMetricsCursorOpen
}

type serverStatusMetricsCursorOpen struct {
	Total  int `metric_name:"cursor.openTotal"  source_type:"gauge"`
	Pinned int `metric_name:"cursor.openPinned" source_type:"gauge"`
}

type serverStatusMetricsCommands struct {
	Count         serverStatusMetricsCommandCount
	CreateIndexes serverStatusMetricsCommandCreateIndexes
	Delete        serverStatusMetricsCommandDelete
	Eval          serverStatusMetricsCommandEval
	FindAndModify serverStatusMetricsCommandFindAndModify
	Insert        serverStatusMetricsCommandInsert
	Update        serverStatusMetricsCommandUpdate
}

type serverStatusMetricsCommandCount struct {
	Failed int `metric_name:"commands.countFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.countFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsCommandCreateIndexes struct {
	Failed int `metric_name:"commands.createIndexesFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.createIndexesFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsCommandDelete struct {
	Failed int `metric_name:"commands.deleteFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.deleteFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsCommandEval struct {
	Failed int `metric_name:"commands.evalFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.evalFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsCommandFindAndModify struct {
	Failed int `metric_name:"commands.modifyFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.modifyFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsCommandInsert struct {
	Failed int `metric_name:"commands.insertFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.insertFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsCommandUpdate struct {
	Failed int `metric_name:"commands.updateFailedPerSecond" source_type:"rate"`
	Total  int `metric_name:"commands.updateFailedPerSecond" source_type:"rate"`
}

type serverStatusMetricsGetLastError struct {
	Wtime     serverStatusMetricsGetLastErrorWtime
	Wtimeouts int `metric_name:"getlasterror.wtimeoutsPerSecond" source_type:"rate"`
}

type serverStatusMetricsGetLastErrorWtime struct {
	TotalMillis int `metric_name:"getlasterror.wtimeMillisPerSecond" source_type:"rate"`
}

type serverStatusMetricsOperation struct {
	ScanAndOrder   int `metric_name:"operation.scanAndOrderPerSecond"   source_type:"rate"`
	WriteConflicts int `metric_name:"operation.writeConflictsPerSecond" source_type:"rate"`
}

type serverStatusMetricsQueryExecutor struct {
	Scanned int `metric_name:"queryexecutor.scannedPerSecond" source_type:"rate"`
}

type serverStatusMetricsRecord struct {
	Moves int `metric_name:"record.movesPerSecond" source_type:"rate"`
}

type serverStatusMetricsRepl struct {
	Apply   serverStatusMetricsReplApply
	Buffer  serverStatusMetricsReplBuffer
	Network serverStatusMetricsReplNetwork
	Preload serverStatusMetricsReplPreload
}

type serverStatusMetricsReplApply struct {
	Ops     int `metric_name:"repl.apply.operationsPerSecond" source_type:"rate"`
	Batches serverStatusMetricsReplApplyBatches
}

type serverStatusMetricsReplApplyBatches struct {
	Num int `metric_name:"repl.apply.batchesPerSec" source_type:"rate"`
}

type serverStatusMetricsReplBuffer struct {
	Count        int `metric_name:"repl.buffer.count"          source_type:"gauge"`
	MaxSizeBytes int `metric_name:"repl.buffer.maxSizeInBytes" source_type:"gauge"`
	SizeBytes    int `metric_name:"repl.buffer.sizeInBytes"    source_type:"gauge"`
}

type serverStatusMetricsReplNetwork struct {
	Bytes          int `metric_name:"repl.network.bytesPerSecond"          source_type:"rate"`
	Ops            int `metric_name:"repl.network.operationPerSecond"      source_type:"rate"`
	ReadersCreated int `metric_name:"repl.network.readersCreatedPerSecond" source_type:"rate"`
	Getmores       serverStatusMetricsReplNetworkGetmores
}

type serverStatusMetricsReplNetworkGetmores struct {
	Num int `metric_name:"repl.network.getmoresPerSecond" source_type:"rate"`
}

type serverStatusMetricsReplPreload struct {
	Docs    serverStatusMetricsReplPreloadDocs
	Indexes serverStatusMetricsReplPreloadIndexes
}

type serverStatusMetricsReplPreloadDocs struct {
	Num         int `metric_name:"repl.docsLoadedPrefetch" source_type:"gauge"`
	TotalMillis int `metric_name:"repl.docsPreloadInMilliseconds" source_type:"gauge"`
}

type serverStatusMetricsReplPreloadIndexes struct {
	Num         int `metric_name:"repl.indexLoadedPrefetch" source_type:"gauge"`
	TotalMillis int `metric_name:"repl.indexPreloadInMilliseconds" source_type:"gauge"`
}

type serverStatusMetricsTtl struct {
	DeletedDocuments int `metric_name:"ttl.deleteDocumentsPerSecond" source_type:"rate"`
	Passes           int `metric_name:"ttl.removeDocumentPerSecond" source_type:"rate"`
}

type serverStatusBackgroundFlushing struct {
	Flushes    int     `metric_name:"flush.flushesDisk"           source_type:"gauge"`
	Total_Ms   float64 `metric_name:"flush.totalInMillisends"     source_type:"gauge"`
	Average_Ms float64 `metric_name:"flush.averageInMilliseconds" source_type:"gauge"`
	Last_Ms    float64 `metric_name:"flush.lastInMilliseconds"    source_type:"gauge"`
}

type serverStatusGlobalLock struct {
	TotalTime     float32 `metric_name:"globallock.totaltime" source_type:"gauge"`
	CurrentQueue  serverStatusGlobalLockCurrentQueue
	ActiveClients serverStatusGlobalLockActiveClients
}

type serverStatusGlobalLockCurrentQueue struct {
	Total   int `metric_name:"globallock.currentQueueTotal"   source_type:"gauge"`
	Readers int `metric_name:"globallock.currentQueueReaders" source_type:"gauge"`
	Writers int `metric_name:"globallock.currentQueueWriters" source_type:"gauge"`
}

type serverStatusGlobalLockActiveClients struct {
	Total   int `metric_name:"globallock.activeClientsTotal"   source_type:"gauge"`
	Readers int `metric_name:"globallock.activeClientsReaders" source_type:"gauge"`
	Writers int `metric_name:"globallock.activeClientsWriters" source_type:"gauge"`
}

type serverStatusExtraInfo struct {
	PageFaults int `metric_name:"pageFaultsPerSecond" source_type:"rate"`
}

type serverStatusWiredTiger struct {
	Cache                  serverStatusWiredTigerCache
	ConcurrentTransactions serverStatusWiredTigerConcurrentTransactions
}

type serverStatusWiredTigerCache struct {
	Size                   int `bson:"bytes currently in the cache"                                 metric_name:"wiredtiger.cacheInBytes"                                 source_type:"gauge"`
	FailedEvictions        int `bson:"failed eviction of pages that exceeded the in-memory maximum" metric_name:"wiredtiger.failedEvictionsPagesPerSecond"                source_type:"gauge"`
	PageSplits             int `bson:"in-memory page splits"                                        metric_name:"cacheInMemoryPageSplits"                                 source_type:"gauge"`
	MaxSize                int `bson:"maximum bytes configured"                                     metric_name:"wiredtiger.cacheMaxInBytes"                              source_type:"gauge"`
	MaxPageSize            int `bson:"maximum page size at eviction"                                metric_name:"wiredtiger.cacheMaxPageSizeEvictionInBytes"              source_type:"gauge"`
	ModifiedPagesEvicted   int `bson:"modified pages evicted"                                       metric_name:"wiredtiger.cacheModifiedPagesEvicted"                    source_type:"gauge"`
	PagesHeld              int `bson:"pages currently held in the cache"                            metric_name:"wiredtiger.cachePagesHeld"                               source_type:"gauge"`
	PagesEvictedThreads    int `bson:"pages evicted by application threads"                         metric_name:"wiredtiger.cachePagesEvictedApplicationThreadsPerSecond" source_type:"gauge"`
	PagesEvictedMax        int `bson:"pages evicted because they exceeded the in-memory maximum"    metric_name:"wiredtiger.cachePagesEvictedInMemoryMaxPerSecond"        source_type:"rate"`
	DirtyData              int `bson:"tracked dirty bytes in the cache"                             metric_name:"wiredtiger.cacheDirtyDataInBytes"                        source_type:"gauge"`
	UnmodifiedPagesEvicted int `bson:"unmodified pages evicted"                                     metric_name:"wiredtiger.cacheUnmodifiedPagesEvicted"                  source_type:"gauge"`
}

type serverStatusWiredTigerConcurrentTransactions struct {
	Write serverStatusWiredTigerConcurrentTransactionsWrite
	Read  serverStatusWiredTigerConcurrentTransactionsRead
}

type serverStatusWiredTigerConcurrentTransactionsWrite struct {
	Out          int `metric_name:"wiredtiger.concurrentTransactions.WriteRemaining" source_type:"gauge"`
	Available    int `metric_name:"wiredtiger.concurrentTransactions.WriteAvailable" source_type:"gauge"`
	TotalTickets int `metric_name:"wiredtiger.concurrentTransactions.WriteTotal"     source_type:"gauge"`
}

type serverStatusWiredTigerConcurrentTransactionsRead struct {
	Out          int `metric_name:"wiredtiger.concurrentTransactions.ReadRemaining" source_type:"gauge"`
	Available    int `metric_name:"wiredtiger.concurrentTransactions.ReadAvailable" source_type:"gauge"`
	TotalTickets int `metric_name:"wiredtiger.concurrentTransactions.ReadTotal"     source_type:"gauge"`
}

type serverStatusLocks struct {
	Collection     serverStatusLocksCollection
	Database       serverStatusLocksDatabase
	Global         serverStatusLocksGlobal
	Metadata       serverStatusLocksMetadata
	GMMAPV1Journal serverStatusLocksMMAPV1Journal
	Oplog          serverStatusLocksOplog
}

type serverStatusLocksCollection struct {
	AcquireCount        serverStatusLocksCollectionAcquireCount
	AcquireWaitCount    serverStatusLocksCollectionAcquireWaitCount
	TimeAcquiringMicros serverStatusLocksCollectionTimeAcquiringMicros
}

type serverStatusLocksCollectionAcquireCount struct {
	Shared          int `bson:"R" metric_name:"locks.collectionAcquiredShared" source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.collectionAcquiredShared" source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.collectionAcquiredShared" source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.collectionAcquiredShared" source_type:"gauge"`
}

type serverStatusLocksCollectionAcquireWaitCount struct {
	Shared    int `bson:"R" metric_name:"locks.collectionAcquireWaitCountShared"    source_type:"gauge"`
	Exclusive int `bson:"W" metric_name:"locks.collectionAcquireWaitCountExclusive" source_type:"gauge"`
}

type serverStatusLocksCollectionTimeAcquiringMicros struct {
	Shared    int `bson:"R" metric_name:"locks.collectionTimeAcquiringMicrosShared"    source_type:"gauge"`
	Exclusive int `bson:"W" metric_name:"locks.collectionTimeAcquiringMicrosExclusive" source_type:"gauge"`
}

type serverStatusLocksDatabase struct {
	AcquireCount        serverStatusLocksDatabaseAcquireCount
	AcquireWaitCount    serverStatusLocksDatabaseAcquireWaitCount
	TimeAcquiringMicros serverStatusLocksDatabaseTimeAcquiringMicros
}

type serverStatusLocksDatabaseAcquireCount struct {
	Shared          int `bson:"R" metric_name:"locks.databaseAcquiredShared" source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.databaseAcquiredShared" source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.databaseAcquiredShared" source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.databaseAcquiredShared" source_type:"gauge"`
}

type serverStatusLocksDatabaseAcquireWaitCount struct {
	Shared          int `bson:"R" metric_name:"locks.databaseAcquireWaitCountShared"          source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.databaseAcquireWaitCountExclusive"       source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.databaseAcquireWaitCountIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.databaseAcquireWaitCountIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksDatabaseTimeAcquiringMicros struct {
	Shared          int `bson:"R" metric_name:"locks.databaseTimeAcquiringMicrosShared"          source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.databaseTimeAcquiringMicrosExclusive"       source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.databaseTimeAcquiringMicrosIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.databaseTimeAcquiringMicrosIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksGlobal struct {
	AcquireCount        serverStatusLocksDatabaseAcquireCount
	AcquireWaitCount    serverStatusLocksDatabaseAcquireWaitCount
	TimeAcquiringMicros serverStatusLocksDatabaseTimeAcquiringMicros
}

type serverStatusLocksGlobalAcquireCount struct {
	Shared          int `bson:"R" metric_name:"locks.globalAcquiredShared" source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.globalAcquiredShared" source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.globalAcquiredShared" source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.globalAcquiredShared" source_type:"gauge"`
}

type serverStatusLocksGlobalAcquireWaitCount struct {
	Shared          int `bson:"R" metric_name:"locks.globalAcquireWaitCountShared"          source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.globalAcquireWaitCountExclusive"       source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.globalAcquireWaitCountIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.globalAcquireWaitCountIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksGlobalTimeAcquiringMicros struct {
	Shared          int `bson:"R" metric_name:"locks.globalTimeAcquiringMicrosShared"          source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.globalTimeAcquiringMicrosExclusive"       source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.globalTimeAcquiringMicrosIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.globalTimeAcquiringMicrosIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksMetadata struct {
	AcquireCount serverStatusLocksMetadataAcquireCount
}

type serverStatusLocksMetadataAcquireCount struct {
	Shared    int `bson:"R" metric_name:"locks.metadataAcquiredShared" source_type:"gauge"`
	Exclusive int `bson:"W" metric_name:"locks.metadataAcquiredShared" source_type:"gauge"`
}

type serverStatusLocksMMAPV1Journal struct {
	AcquireCount        serverStatusLocksMMAPV1JournalAcquireCount
	TimeAcquiringMicros serverStatusLocksMMAPV1JournalTimeAcquiringMicros
}

type serverStatusLocksMMAPV1JournalAcquireCount struct {
	Shared          int `bson:"R" metric_name:"locks.mmapv1journalAcquiredShared"          source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.mmapv1journalAcquiredExclusive"       source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.mmapv1journalAcquiredIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.mmapv1journalAcquiredIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksOplog struct {
	AcquireCount        serverStatusLocksOplogAcquireCount
	TimeAcquiringMicros serverStatusLocksOplogTimeAcquiringMicros
}

type serverStatusLocksOplogAcquireCount struct {
	Shared          int `bson:"R" metric_name:"locks.oplogAcquiredShared"          source_type:"gauge"`
	Exclusive       int `bson:"W" metric_name:"locks.oplogAcquiredExclusive"       source_type:"gauge"`
	IntentShared    int `bson:"r" metric_name:"locks.oplogAcquiredIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.oplogAcquiredIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksOplogTimeAcquiringMicros struct {
	IntentShared    int `bson:"r" metric_name:"locks.oplogTimeAcquiringMicrosIntentShared"    source_type:"gauge"`
	IntentExclusive int `bson:"w" metric_name:"locks.oplogTimeAcquiringMicrosIntentExclusive" source_type:"gauge"`
}

type serverStatusLocksMMAPV1JournalTimeAcquiringMicros struct {
	Shared    int `bson:"R" metric_name:"locks.mmapv1journalTimeAcquiringMicrosShared"          source_type:"gauge"`
	Exclusive int `bson:"W" metric_name:"locks.mmapv1journalTimeAcquiringMicrosExclusive"       source_type:"gauge"`
}

type serverStatusDur struct {
	Commits               int     `metric_name:"dur.commits"                            source_type:"gauge"`
	CommitsInWriteLock    int     `metric_name:"dur.commitsInWriteLock"                 source_type:"gauge"`
	Compression           float32 `metric_name:"dur.compression"                        source_type:"gauge"`
	EarlyCommits          int     `metric_name:"dur.earlyCommits"                       source_type:"gauge"`
	JournaledMB           float32
	JournaledBytes        float32 `metric_name:"dur.journaledInBytes"                   source_type:"gauge"`
	WriteToDataFilesMB    float32
	WriteToDataFilesBytes float32 `metric_name:"dur.dataWrittenJournalDataFilesInBytes" source_type:"gauge"`
	TimeMS                serverStatusDurTimeMS
}

type serverStatusDurTimeMS struct {
	Commits            int `metric_name:"dur.commitsInMilliseconds"              source_type:"gauge"`
	CommitsInWriteLock int `metric_name:"dur.commitsInWriteLockInMilliseconds"   source_type:"gauge"`
	Dt                 int `metric_name:"dur.timeCollectedCommitsInMilliseconds" source_type:"gauge"`
	PrepLogBuffer      int `metric_name:"dur.preparingInMilliseconds"            source_type:"gauge"`
	RemapPrivateView   int `metric_name:"dur.remappingInMilliseconds"            source_type:"gauge"`
	WriteToDataFiles   int `metric_name:"dur.writingDataFilesInMilliseconds"     source_type:"gauge"`
	WriteToJournal     int `metric_name:"dur.writingJournalInMilliseconds"       source_type:"gauge"`
}

type collStats struct {
	Size        int  `metric_name:"collection.sizeInBytes"        source_type:"gauge"`
	AvgObjSize  int  `metric_name:"collection.avgObjSizeInBytes"  source_type:"gauge"`
	Count       int  `metric_name:"collection.count"              source_type:"gauge"`
	Capped      bool `metric_name:"collection.capped"             source_type:"attribute"`
	Max         int  `metric_name:"collection.max"                source_type:"gauge"`
	MaxSize     int  `metric_name:"collection.maxSizeInBytes"     source_type:"gauge"`
	StorageSize int  `metric_name:"collection.storageSizeInBytes" source_type:"gauge"`
	Nindexes    int  `metric_name:"collection.nindexes"           source_type:"gauge"`
}

type dbStats struct {
	Objects     int `metric_name:"stats.objects"        source_type:"gauge"`
	StorageSize int `metric_name:"stats.storageInBytes" source_type:"gauge"`
	IndexSize   int `metric_name:"stats.indexInBytes"   source_type:"gauge"`
	Indexes     int `metric_name:"stats.indexes"        source_type:"gauge"`
	DataSize    int `metric_name:"stats.dataInBytes"    source_type:"gauge"`
}

type replSetGetStatus struct {
	Members []replSetGetStatusMember
}

type replSetGetStatusMember struct {
	Name     string `metric_name:"replset.name" source_type:"attribute"`
	Health   int    `metric_name:"replset.health" source_type:"gauge"`
	StateStr string `metric_name:"replset.state" source_type:"attribute"`
	Uptime   int    `metric_name:"replset.uptimeInMilliseconds" source_type:"gauge"`
}

package metrics

// CollStats is a storage struct for unmarshalling the collStats command
type CollStats struct {
	Size        *int            `bson:"size"        metric_name:"collection.sizeInBytes"        source_type:"gauge"`
	AvgObjSize  *int            `bson:"avgObjSize"  metric_name:"collection.avgObjSizeInBytes"  source_type:"gauge"`
	Count       *int            `bson:"count"       metric_name:"collection.count"              source_type:"gauge"`
	Capped      *bool           `bson:"capped"      metric_name:"collection.capped"             source_type:"gauge"`
	Max         *int            `bson:"max"         metric_name:"collection.max"                source_type:"gauge"`
	MaxSize     *int            `bson:"maxSize"     metric_name:"collection.maxSizeInBytes"     source_type:"gauge"`
	StorageSize *int            `bson:"storageSize" metric_name:"collection.storageSizeInBytes" source_type:"gauge"`
	Nindexes    *int            `bson:"nindexes"    metric_name:"collection.nindexes"           source_type:"gauge"`
	IndexSizes  *map[string]int `bson:"indexSizes"`
}

// IndexStats is a storage struct for $indexStats
type IndexStats struct {
	Name     *string             `bson:"name"`
	Accesses *IndexStatsAccesses `bson:"accesses"`
}

// IndexStatsAccesses is a storage struct for $indexStats
type IndexStatsAccesses struct {
	Ops *int `bson:"ops"`
}

package metrics

type DbStats struct {
	Objects     *int `bson:"objects"     metric_name:"stats.objects"        source_type:"gauge"`
	StorageSize *int `bson:"storageSize" metric_name:"stats.storageInBytes" source_type:"gauge"`
	IndexSize   *int `bson:"indexSize"   metric_name:"stats.indexInBytes"   source_type:"gauge"`
	Indexes     *int `bson:"indexes"     metric_name:"stats.indexes"        source_type:"gauge"`
	DataSize    *int `bson:"dataSize"    metric_name:"stats.dataInBytes"    source_type:"gauge"`
}

package metrics

// ReplSetGetStatus is a storage struct
type ReplSetGetStatus struct {
	Members []ReplSetGetStatusMember `bson:"members"`
}

// ReplSetGetStatusMember is a storage struct
type ReplSetGetStatusMember struct {
	Name     *string `bson:"name"`
	Health   *int    `bson:"health"   metric_name:"replset.health"               source_type:"gauge"`
	StateStr *string `bson:"stateStr" metric_name:"replset.state"                source_type:"attribute"`
	Uptime   *int    `bson:"uptime"   metric_name:"replset.uptimeInMilliseconds" source_type:"gauge"`
}

// IsMaster is a storage struct
type IsMaster struct {
	SetName   *string `bson:"setName"`
	IsMaster  *bool   `bson:"ismaster"  metric_name:"replset.isMaster"    source_type:"gauge"`
	Secondary *bool   `bson:"secondary" metric_name:"replset.isSecondary" source_type:"gauge"`
}

// ReplSetGetConfig is a storage struct
type ReplSetGetConfig struct {
	Config struct {
		Members []struct {
			Host        *string  `bson:"host"`
			ArbiterOnly *bool    `bson:"arbiterOnly" metric_name:"replset.isArbiter" source_type:"gauge"`
			Hidden      *bool    `bson:"hidden"      metric_name:"replset.isHidden"  source_type:"gauge"`
			Priority    *float32 `bson:"priority"    metric_name:"replset.priority"  source_type:"gauge"`
			Votes       *float32 `bson:"votes"       metric_name:"replset.votes"     source_type:"gauge"`
		} `bson:"members"`
	} `bson:"config"`
}

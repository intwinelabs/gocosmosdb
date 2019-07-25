package gocosmosdb

// Resource
type Resource struct {
	Id   string `json:"id,omitempty"`
	Self string `json:"_self,omitempty"`
	Etag string `json:"_etag,omitempty"`
	Rid  string `json:"_rid,omitempty"`
	Ts   int    `json:"_ts,omitempty"`
}

// Indexing policy
type IndexingPolicy struct {
	Automatic     bool `json:"automatic,omitempty"`
	IncludedPaths []struct {
		Indexes []struct {
			DataType  string `json:"dataType,omitempty"`
			Kind      string `json:"kind,omitempty"`
			Precision int    `json:"precision,omitempty`
		} `json:"indexes,omitempty"`
		Path string `json:"path,omitempty"`
	} `json:"includedPaths,omitempty"`
	IndexingMode string `json:"indexingMode,omitempty"`
}

// Partition Key
type PartitionKeyDef struct {
	Kind  string   `json:"kind"`
	Paths []string `json:"paths"`
}

// Database
type Database struct {
	Resource
	Colls string `json:"_colls,omitempty"`
	Users string `json:"_users,omitempty"`
}

// Collection
type Collection struct {
	Resource
	IndexingPolicy  IndexingPolicy  `json:"indexingPolicy,omitempty"`
	PartitionKeyDef PartitionKeyDef `json:"partitionKey,omitempty"`
	Docs            string          `json:"_docs,omitempty"`
	Udf             string          `json:"_udfs,omitempty"`
	Sporcs          string          `json:"_sporcs,omitempty"`
	Triggers        string          `json:"_triggers,omitempty"`
	Conflicts       string          `json:"_conflicts,omitempty"`
}

// QueryWithParameters
type QueryWithParameters struct {
	Query      string          `json:"query"`
	Parameters []QueryParamter `json:"parameters"`
}

// QueryParameter
type QueryParamter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Document
type Document struct {
	Resource
	Attachments string `json:"attachments,omitempty"`
}

// Expirable
type Expirable struct {
	TTL int64 `json:"ttl"`
}

// Stored Procedure
type Sproc struct {
	Resource
	Body string `json:"body,omitempty"`
}

// User Defined Function
type UDF struct {
	Resource
	Body string `json:"body,omitempty"`
}

// Metrics
type Metrics struct {
	TotalExecutionTimeInMs         float64 `json:"totalExecutionTimeInMs,omitempty"`
	QueryCompileTimeInMs           float64 `json:"queryCompileTimeInMs,omitempty"`
	QueryLogicalPlanBuildTimeInMs  float64 `json:"queryLogicalPlanBuildTimeInMs,omitempty"`
	QueryPhysicalPlanBuildTimeInMs float64 `json:"queryPhysicalPlanBuildTimeInMs,omitempty"`
	QueryOptimizationTimeInMs      float64 `json:"queryOptimizationTimeInMs,omitempty"`
	VMExecutionTimeInMs            float64 `json:"VMExecutionTimeInMs,omitempty"`
	IndexLookupTimeInMs            float64 `json:"indexLookupTimeInMs,omitempty"`
	DocumentLoadTimeInMs           float64 `json:"documentLoadTimeInMs,omitempty"`
	SystemFunctionExecuteTimeInMs  float64 `json:"systemFunctionExecuteTimeInMs,omitempty"`
	UserFunctionExecuteTimeInMs    float64 `json:"userFunctionExecuteTimeInMs,omitempty"`
	RetrievedDocumentCount         int     `json:"retrievedDocumentCount,omitempty"`
	RetrievedDocumentSize          int     `json:"retrievedDocumentSize,omitempty"`
	OutputDocumentCount            int     `json:"outputDocumentCount,omitempty"`
	WriteOutputTimeInMs            float64 `json:"writeOutputTimeInMs,omitempty"`
	IndexUtilizationRatio          float64 `json:"indexUtilizationRatio,omitempty"`
	RequestCharge                  float64 `json:"requestCharge,omitempty"`
}

// PartitionKeyRange partition key range model
type PartitionKeyRange struct {
	Resource
	PartitionKeyRangeID string `json:"id,omitempty"`
	MinInclusive        string `json:"minInclusive,omitempty"`
	MaxInclusive        string `json:"maxExclusive,omitempty"`
}

// PagableQuery
type PagableQuery struct {
	client       *CosmosDB
	coll         string
	query        *QueryWithParameters
	sessionToken CallOption
	continuation CallOption
	limit        CallOption
	offset       int64
	docs         interface{}
}

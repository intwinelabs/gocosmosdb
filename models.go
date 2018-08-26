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
// TODO: Ex/IncludePaths
type IndexingPolicy struct {
	IndexingMode string `json:"indexingMode,omitempty"`
	Automatic    bool   `json:"automatic,omitempty"`
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
	IndexingPolicy IndexingPolicy `json:"indexingPolicy,omitempty"`
	Docs           string         `json:"_docs,omitempty"`
	Udf            string         `json:"_udfs,omitempty"`
	Sporcs         string         `json:"_sporcs,omitempty"`
	Triggers       string         `json:"_triggers,omitempty"`
	Conflicts      string         `json:"_conflicts,omitempty"`
}

// Document
type Document struct {
	Resource
	Attachments string `json:"attachments,omitempty"`
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

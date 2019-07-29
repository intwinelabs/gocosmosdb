package gocosmosdb

const (
	// HeaderActivityID - A client supplied identifier for the operation, which is echoed in the server response.
	// The recommended value is a unique identifier.
	HeaderActivityID = "X-Ms-Activity-Id"

	// HeaderAIM - Indicates a change feed request. Must be set to "Incremental feed", or omitted otherwise.
	HeaderAIM = "A-IM"

	//HeaderAllowTenativeWrites - For using multiple write locations.
	HeaderAllowTenativeWrites = "X-Ms-Cosmos-Allow-Tentative-Writes"

	// HeaderAuth - The authorization token for the request
	HeaderAuth = "Authorization"

	// HeaderConsistencyLevel - The consistency level override for read options against documents and attachments.
	// The valid values are: Strong, Bounded, Session, or Eventual
	HeaderConsistencyLevel = "X-Ms-Consistency-Level"

	// HeaderContentLength - Indicates the size of the entity-body, in bytes, sent to the recipient.
	HeaderContentLength = "Content-Length"

	// HeaderContentType - POST it must be application/query+json
	// attachments must be set to the Mime type of the attachment
	// all other tasks must be application/json
	HeaderContentType = "Content-Type"

	// HeaderContinuation - A string token returned for queries and read-feed operations if there are more
	// results to be read. Clients can retrieve the next page of results by resubmitting the request with this value.
	HeaderContinuation = "X-Ms-Continuation"

	// HeaderCrossPartition - When this header is set to true and if your query doesn't have a partition key, Azure
	// Cosmos DB fans out the query across partitions. The fan out is done by issuing individual queries to all the
	// partitions. To read the query results, the client applications should consume the results from the FeedResponse
	// and check for the ContinuationToken property. To read all the results, keep iterating on the data until the
	// ContinuationToken is null.
	HeaderCrossPartition = "X-Ms-Documentdb-Query-Enablecrosspartition"

	// HeaderEnableScan - Use an index scan to process the query if the right index path of type is not available.
	HeaderEnableScan = "X-Ms-Documentdb-Query-Enable-Scan"

	// HeaderIfMatch - Used to make operation conditional for optimistic concurrency.
	// The value should be the etag value of the resource.
	HeaderIfMatch = "If-Match"

	// HeaderIfModifiedSince - Returns etag of resource modified after specified date in RFC 1123 format.
	// Ignored when If-None-Match is specified
	HeaderIfModifiedSince = "If-Modified-Since"

	// HeaderIfNonMatch - Makes operation conditional to only execute if the resource has changed.
	// The value should be the etag of the resource.
	HeaderIfNonMatch = "If-None-Match"

	// HeaderIndexingDirective - Overide the collections default indexing policy, set to Include or Exclude.
	HeaderIndexingDirective = "x-ms-indexing-directive"

	// HeaderIsQuery - Required for queries. This property must be set to true.
	HeaderIsQuery = "X-Ms-Documentdb-Isquery"

	// HeaderIsQueryPlan -
	HeaderIsQueryPlan = "X-Ms-Cosmos-Is-Query-Plan-Request"

	// HeaderMaxItemCount - An integer indicating the maximum number of items to be returned per page.
	// An x-ms-max-item-count of -1 can be specified to let the service determine the optimal item count.
	HeaderMaxItemCount = "X-Ms-Max-Item-Count"

	// HeaderOfferThroughput - The user specified throughput for the collection expressed in units of 100
	// request units per second.
	HeaderOfferThroughput = "X-Ms-Offer-Throughput"

	// HeaderParalelizeCrossPartition - Sets the query to run in parallel across partitions.
	HeaderParalelizeCrossPartition = "X-Ms-Documentdb-Query-Parallelizecrosspartitionquery"

	// HeaderPartitionKey - The partition key value for the requested document or attachment operation.
	// Required for operations against documents and attachments when the collection definition includes
	// a partition key definition. This value is used to scope your query to documents that match the partition
	// key criteria. By design it's a single partition query. Supported in API versions 2015-12-16 and newer.
	// Currently, the SQL API supports a single partition key, so this is an array containing just one value.
	HeaderPartitionKey = "X-Ms-Documentdb-Partitionkey"

	// HeaderPartitionKeyRangeID - Used in change feed requests. The partition key range ID for reading data.
	HeaderPartitionKeyRangeID = "X-Ms-Documentdb-Partitionkeyrangeid"

	// HeaderPopulateQueryMetrics - Set to obtain detailed metrics on query execution.
	HeaderPopulateQueryMetrics = "X-Ms-Documentdb-Populatequerymetrics"

	// HeaderQueryMetrics - The query statistics for the execution. This is a delimited string containing statistics
	// of time spent in the various phases of query execution.
	HeaderQueryMetrics = "X-Ms-Documentdb-Query-Metrics"

	// HeaderQueryVersion - Set the query version.
	HeaderQueryVersion = "X-Ms-Cosmos-Query-Version"

	// HeaderRequestCharge - The number of request units consumed by the operation.
	HeaderRequestCharge = "X-Ms-Request-Charge"

	// HeaderSessionToken - A string token used with session level consistency.
	HeaderSessionToken = "X-Ms-Session-Token"

	// HeaderSupportedQueryFeatures -
	HeaderSupportedQueryFeatures = "X-Ms-Cosmos-Supported-Query-Features"

	// HeaderUpsert - If set to true, Cosmos DB creates the document with the ID (and partition key value if applicable)
	// if it doesn’t exist, or update the document if it exists.
	HeaderUpsert = "X-Ms-Documentdb-Is-Upsert"

	// HeaderUserAgent - A string that specifies the client user agent performing the request.
	// The recommended format is {user agent name}/{version}.
	HeaderUserAgent = "User-Agent"

	// HeaderVersion - The version of the Cosmos DB REST service.
	HeaderVersion = "X-Ms-Version"

	// HeaderXDate - The date of the request per RFC 1123 date format expressed in Coordinated Universal Time.
	// For example, Fri, 08 Apr 2015 03:52:31 GMT.
	HeaderXDate = "X-Ms-Date"
)

package gocosmosdb

import (
	"context"
	"encoding/json"
	"strconv"
)

// Consistency type to define consistency levels
type Consistency string

const (
	// Strong consistency level
	Strong Consistency = "strong"

	// Bounded consistency level
	Bounded Consistency = "bounded"

	// Session consistency level
	Session Consistency = "session"

	// Eventual consistency level
	Eventual Consistency = "eventual"
)

// CallOption function
type CallOption func(r *Request) error

// PartitionKey - specificy which partiotion will be used to satisfty the request
func PartitionKey(partitionKey interface{}) CallOption {

	// The partition key header must be an array following the spec:
	// https: //docs.microsoft.com/en-us/rest/api/cosmos-db/common-cosmosdb-rest-request-headers
	// and must contain brackets
	// example: x-ms-documentdb-partitionkey: [ "abc" ]
	var (
		pk  []byte
		err error
	)
	switch v := partitionKey.(type) {
	case json.Marshaler:
		pk, err = json.Marshal(v)
	default:
		pk, err = json.Marshal([]interface{}{v})
	}

	header := []string{string(pk)}

	return func(r *Request) error {
		if err != nil {
			return err
		}
		r.Header[HeaderPartitionKey] = header
		return nil
	}
}

// Upsert - if set to true, Cosmos DB creates the document with the ID (and partition key value if applicable) if it doesn’t exist, or update the document if it exists.
func Upsert() CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderUpsert, "true")
		return nil
	}
}

// Limit - set max item count for response
func Limit(limit int) CallOption {
	header := strconv.Itoa(limit)
	return func(r *Request) error {
		r.Header.Set(HeaderMaxItemCount, header)
		return nil
	}
}

// Continuation - a string token returned for queries and read-feed operations if there are more results to be read. Clients can retrieve the next page of results by resubmitting the request with the x-ms-continuation request header set to this value.
func Continuation(continuation string) CallOption {
	return func(r *Request) error {
		if continuation == "" {
			return nil
		}
		r.Header.Set(HeaderContinuation, continuation)
		return nil
	}
}

// ConsistencyLevel - override for read options against documents and attachments. The valid values are: Strong, Bounded, Session, or Eventual (in order of strongest to weakest). The override must be the same or weaker than the account�s configured consistency level.
func ConsistencyLevel(consistency Consistency) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderConsistencyLevel, string(consistency))
		return nil
	}
}

// SessionToken - a string token used with session level consistency.
func SessionToken(sessionToken string) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderSessionToken, sessionToken)
		return nil
	}
}

// CrossPartition - allows query to run on all partitions
func CrossPartition() CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderCrossPartition, "true")
		return nil
	}
}

// IfMatch - used to make operation conditional for optimistic concurrency. The value should be the etag value of the resource.
// (applicable only on PUT and DELETE)
func IfMatch(eTag string) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderIfMatch, eTag)
		return nil
	}
}

// IfNoneMatch - makes operation conditional to only execute if the resource has changed. The value should be the etag of the resource.
// Optional (applicable only on GET)
func IfNoneMatch(eTag string) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderIfNonMatch, eTag)
		return nil
	}
}

// IfModifiedSince - returns etag of resource modified after specified date in RFC 1123 format. Ignored when If-None-Match is specified
// Optional (applicable only on GET)
func IfModifiedSince(date string) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderIfModifiedSince, date)
		return nil
	}
}

// ChangeFeed - indicates a change feed request
func ChangeFeed() CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderAIM, "Incremental feed")
		return nil
	}
}

// ThroughputRUs - adds throughput headers for container creation
func ThroughputRUs(rus int) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderOfferThroughput, strconv.Itoa(rus))
		return nil
	}
}

// PartitionKeyRangeID - adds the partition key range header
func PartitionKeyRangeID(id int) CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderPartitionKeyRangeID, strconv.Itoa(id))
		return nil
	}
}

// EnableQueryScan - add the scan header
func EnableQueryScan() CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderEnableScan, "true")
		return nil
	}
}

// EnableParallelizeCrossPartitionQuery - add the parallelize header
func EnableParallelizeCrossPartitionQuery() CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderParalelizeCrossPartition, "true")
		return nil
	}
}

// EnablePopulateQueryMetrics - add the parallelize header
func EnablePopulateQueryMetrics() CallOption {
	return func(r *Request) error {
		r.Header.Set(HeaderPopulateQueryMetrics, "true")
		return nil
	}
}

// WithContext - adds a context to the request
func WithContext(ctx context.Context) CallOption {
	return func(r *Request) error {
		r.rContext = ctx
		return nil
	}
}

// QueryVersion - add the query latest query version header
func QueryVersion() CallOption {
	return func(r *Request) error {
		r.Header.Add(HeaderQueryVersion, SupportedQueryVersion)
		return nil
	}
}

package gocosmosdb

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallOptions(t *testing.T) {
	assert := assert.New(t)
	opts := []CallOption{}
	opts = append(opts, PartitionKey("test"))
	opts = append(opts, Upsert())
	opts = append(opts, Limit(100))
	opts = append(opts, Continuation("continueToken"))
	opts = append(opts, ConsistencyLevel("Strong"))
	opts = append(opts, SessionToken("sessionToken"))
	opts = append(opts, CrossPartition())
	opts = append(opts, IfMatch("eTag"))
	opts = append(opts, IfNoneMatch("eTag"))
	opts = append(opts, IfModifiedSince("someDate"))
	opts = append(opts, ChangeFeed())
	opts = append(opts, ThroughputRUs(400))
	opts = append(opts, PartitionKeyRangeID(0))
	opts = append(opts, EnableQueryScan())
	opts = append(opts, EnableParallelizeCrossPartitionQuery())
	opts = append(opts, EnablePopulateQueryMetrics())
	ctx := context.WithValue(context.Background(), "foo", "bar")
	opts = append(opts, WithContext(ctx))

	link := "http://localhost:8080"
	req, err := http.NewRequest("POST", link, nil)
	if err != nil {
		assert.Nil(err)
	}
	r := ResourceRequest(link, req)
	for i := 0; i < len(opts); i++ {
		if err := opts[i](r); err != nil {
			assert.Nil(err)
		}
	}

	assert.Equal("[\"test\"]", r.Header.Get(HeaderPartitionKey))
	assert.Equal("true", r.Header.Get(HeaderUpsert))
	assert.Equal("100", r.Header.Get(HeaderMaxItemCount))
	assert.Equal("continueToken", r.Header.Get(HeaderContinuation))
	assert.Equal("Strong", r.Header.Get(HeaderConsistencyLevel))
	assert.Equal("sessionToken", r.Header.Get(HeaderSessionToken))
	assert.Equal("true", r.Header.Get(HeaderCrossPartition))
	assert.Equal("eTag", r.Header.Get(HeaderIfMatch))
	assert.Equal("eTag", r.Header.Get(HeaderIfNonMatch))
	assert.Equal("someDate", r.Header.Get(HeaderIfModifiedSince))
	assert.Equal("Incremental feed", r.Header.Get(HeaderAIM))
	assert.Equal("400", r.Header.Get(HeaderOfferThroughput))
	assert.Equal("0", r.Header.Get(HeaderPartitionKeyRangeID))
	assert.Equal("true", r.Header.Get(HeaderEnableScan))
	assert.Equal("true", r.Header.Get(HeaderParalelizeCrossPartition))
	assert.Equal("true", r.Header.Get(HeaderPopulateQueryMetrics))
	assert.Equal(ctx, r.rContext)
}

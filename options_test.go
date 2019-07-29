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

	assert.Equal([]string{"[\"test\"]"}, r.Header[HeaderPartitionKey])
	assert.Equal([]string{"true"}, r.Header[HeaderUpsert])
	assert.Equal([]string{"100"}, r.Header[HeaderMaxItemCount])
	assert.Equal([]string{"continueToken"}, r.Header[HeaderContinuation])
	assert.Equal([]string{"Strong"}, r.Header[HeaderConsistencyLevel])
	assert.Equal([]string{"sessionToken"}, r.Header[HeaderSessionToken])
	assert.Equal([]string{"true"}, r.Header[HeaderCrossPartition])
	assert.Equal([]string{"eTag"}, r.Header[HeaderIfMatch])
	assert.Equal([]string{"eTag"}, r.Header[HeaderIfNonMatch])
	assert.Equal([]string{"someDate"}, r.Header[HeaderIfModifiedSince])
	//assert.Equal([]string{"Incremental feed"}, r.Header[HeaderAIM])
	assert.Equal([]string{"400"}, r.Header[HeaderOfferThroughput])
	assert.Equal([]string{"0"}, r.Header[HeaderPartitionKeyRangeID])
	assert.Equal([]string{"true"}, r.Header[HeaderEnableScan])
	assert.Equal([]string{"true"}, r.Header[HeaderParalelizeCrossPartition])
	assert.Equal(ctx, r.rContext)
}

package gocosmosdb

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpectStatusCodes(t *testing.T) {
	assert := assert.New(t)

	st1 := expectStatusCode(200)
	assert.Equal(true, st1(200))

	st2 := expectStatusCodeXX(400)
	assert.Equal(true, st2(499))

}

func TestResponseContinuation(t *testing.T) {
	assert := assert.New(t)

	resp := &Response{Header: http.Header{}}
	resp.Header.Set(HeaderContinuation, "testContinuation")
	continuation := resp.Continuation()
	assert.Equal("testContinuation", continuation)
}

func TestResponseSessionToken(t *testing.T) {
	assert := assert.New(t)

	resp := &Response{Header: http.Header{}}
	resp.Header.Set(HeaderSessionToken, "testSession")
	session := resp.SessionToken()
	assert.Equal("testSession", session)
}
func TestGetQueryMetrics(t *testing.T) {
	assert := assert.New(t)

	xMsDocumentdbQueryMetrics := "totalExecutionTimeInMs=33.67;queryCompileTimeInMs=0.06;queryLogicalPlanBuildTimeInMs=0.02;queryPhysicalPlanBuildTimeInMs=0.10;queryOptimizationTimeInMs=0.00;VMExecutionTimeInMs=32.56;indexLookupTimeInMs=0.36;documentLoadTimeInMs=9.58;systemFunctionExecuteTimeInMs=0.00;userFunctionExecuteTimeInMs=0.00;retrievedDocumentCount=2000;retrievedDocumentSize=1125600;outputDocumentCount=2000;writeOutputTimeInMs=18.10;indexUtilizationRatio=1.00"
	xMsRequestCharge := "604.42"
	resp := &Response{Header: http.Header{}}
	resp.Header.Set(HeaderQueryMetrics, xMsDocumentdbQueryMetrics)
	resp.Header.Set(HeaderRequestCharge, xMsRequestCharge)
	metrics, err := resp.GetQueryMetrics()
	_metrics := &Metrics{
		TotalExecutionTimeInMs:         33.67,
		QueryCompileTimeInMs:           0.06,
		QueryLogicalPlanBuildTimeInMs:  0.02,
		QueryPhysicalPlanBuildTimeInMs: 0.1,
		QueryOptimizationTimeInMs:      0,
		VMExecutionTimeInMs:            32.56,
		IndexLookupTimeInMs:            0.36,
		DocumentLoadTimeInMs:           9.58,
		SystemFunctionExecuteTimeInMs:  0,
		UserFunctionExecuteTimeInMs:    0,
		RetrievedDocumentCount:         2000,
		RetrievedDocumentSize:          1125600,
		OutputDocumentCount:            2000,
		WriteOutputTimeInMs:            18.1,
		IndexUtilizationRatio:          1,
		RequestCharge:                  604.42,
	}
	assert.Nil(err)
	assert.Equal(_metrics, metrics)
}

func TestGetQueryMetricsEmpty(t *testing.T) {
	assert := assert.New(t)

	xMsDocumentdbQueryMetricsEmpty := ""
	xMsRequestCharge := "604.42"
	resp := &Response{Header: http.Header{}}
	resp.Header.Set(HeaderQueryMetrics, xMsDocumentdbQueryMetricsEmpty)
	resp.Header.Set(HeaderRequestCharge, xMsRequestCharge)
	metrics, err := resp.GetQueryMetrics()
	assert.Contains(err.Error(), "no metrics in response")
	assert.Nil(metrics)
}

func TestGetRUs(t *testing.T) {
	assert := assert.New(t)

	xMsRequestCharge := "604.42"
	resp := &Response{Header: http.Header{}}
	resp.Header.Set(HeaderRequestCharge, xMsRequestCharge)
	rus, err := resp.GetRUs()
	assert.Nil(err)
	assert.Equal(float64(604.42), rus)
}

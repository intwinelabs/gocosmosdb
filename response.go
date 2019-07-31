package gocosmosdb

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

type Response struct {
	Header http.Header
}

// Continuation - returns continuation token for paged request.
// Pass this value to next request to get next page of documents.
func (r *Response) Continuation() string {
	return r.Header.Get(HeaderContinuation)
}

// SessionToken - returns session token for session consistent request.
// Pass this value to next request to maintain session consistency documents.
func (r *Response) SessionToken() string {
	return r.Header.Get(HeaderSessionToken)
}

// GetResponseMetrics - returns a responses metrics
func (r *Response) GetResponseMetrics() (*Metrics, error) {
	// x-ms-documentdb-query-metrics: totalExecutionTimeInMs=33.67;queryCompileTimeInMs=0.06;queryLogicalPlanBuildTimeInMs=0.02;queryPhysicalPlanBuildTimeInMs=0.10;queryOptimizationTimeInMs=0.00;VMExecutionTimeInMs=32.56;indexLookupTimeInMs=0.36;documentLoadTimeInMs=9.58;systemFunctionExecuteTimeInMs=0.00;userFunctionExecuteTimeInMs=0.00;retrievedDocumentCount=2000;retrievedDocumentSize=1125600;outputDocumentCount=2000;writeOutputTimeInMs=18.10;indexUtilizationRatio=1.00
	// x-ms-request-charge: 604.42

	metrics := &Metrics{}
	metricsStrSlice := strings.Split(r.Header.Get(HeaderQueryMetrics), ";")
	for _, metricStr := range metricsStrSlice {
		metricSlice := strings.Split(metricStr, "=")
		if len(metricSlice) == 2 {
			metricKey, metricVal, err := getMetricKeyVal(metricSlice)
			if err != nil {
				return nil, err
			}
			switch metricKey {
			case "totalExecutionTimeInMs":
				metrics.TotalExecutionTimeInMs = metricVal
			case "queryCompileTimeInMs":
				metrics.QueryCompileTimeInMs = metricVal
			case "queryLogicalPlanBuildTimeInMs":
				metrics.QueryLogicalPlanBuildTimeInMs = metricVal
			case "queryPhysicalPlanBuildTimeInMs":
				metrics.QueryPhysicalPlanBuildTimeInMs = metricVal
			case "queryOptimizationTimeInMs":
				metrics.QueryOptimizationTimeInMs = metricVal
			case "VMExecutionTimeInMs":
				metrics.VMExecutionTimeInMs = metricVal
			case "indexLookupTimeInMs":
				metrics.IndexLookupTimeInMs = metricVal
			case "documentLoadTimeInMs":
				metrics.DocumentLoadTimeInMs = metricVal
			case "systemFunctionExecuteTimeInMs":
				metrics.SystemFunctionExecuteTimeInMs = metricVal
			case "userFunctionExecuteTimeInMs":
				metrics.UserFunctionExecuteTimeInMs = metricVal
			case "retrievedDocumentCount":
				metrics.RetrievedDocumentCount = int(metricVal)
			case "retrievedDocumentSize":
				metrics.RetrievedDocumentSize = int(metricVal)
			case "outputDocumentCount":
				metrics.OutputDocumentCount = int(metricVal)
			case "writeOutputTimeInMs":
				metrics.WriteOutputTimeInMs = metricVal
			case "indexUtilizationRatio":
				metrics.IndexUtilizationRatio = metricVal
			}

		}
	}

	var requestVal float64
	_, err := fmt.Sscanf(r.Header.Get(HeaderRequestCharge), "%f", &requestVal)
	if err != nil {
		return nil, fmt.Errorf("error parsing request charge header: %v", err)

	}
	metrics.RequestCharge = requestVal

	return metrics, nil
}

func getMetricKeyVal(metricSlice []string) (string, float64, error) {
	if len(metricSlice) == 2 {
		var metric float64
		_, err := fmt.Sscanf(metricSlice[1], "%f", &metric)
		if err != nil {
			return "", -1, fmt.Errorf("error parsing metrics header: %v: %s", metricSlice, err)
		}
		return metricSlice[0], metric, nil
	}
	return "", -1, fmt.Errorf("pass metric slice must have a length of 2")
}

type statusCodeValidatorFunc func(statusCode int) bool

func expectStatusCode(expected int) statusCodeValidatorFunc {
	return func(statusCode int) bool {
		return expected == statusCode
	}
}

func expectStatusCodeXX(expected int) statusCodeValidatorFunc {
	begining := int(math.Floor(float64(expected/100))) * 100
	end := begining + 99
	return func(statusCode int) bool {
		return (statusCode >= begining) && (statusCode <= end)
	}
}

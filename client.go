package gocosmosdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/intwinelabs/logger"
	"github.com/moul/http2curl"
)

// Client - struct to hold the underlying SQL REST API client
type apiClient struct {
	uri        string
	config     Config
	httpClient *retryablehttp.Client
	logger     *logger.Logger
}

func newAPIClient(conf *Config) *apiClient {
	client := &apiClient{}
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	client.httpClient = httpClient
	var zeroDuration time.Duration
	if conf.RetryWaitMin == zeroDuration {
		client.httpClient.RetryWaitMin = 10 * time.Millisecond
	} else {
		client.httpClient.RetryWaitMin = conf.RetryWaitMin
	}
	if conf.RetryWaitMax == zeroDuration {
		client.httpClient.RetryWaitMax = 50 * time.Millisecond
	} else {
		client.httpClient.RetryWaitMax = conf.RetryWaitMax
	}
	client.httpClient.RetryMax = conf.RetryMax
	if conf.Pooled {
		client.httpClient.HTTPClient.Transport = cleanhttp.DefaultPooledTransport()
	}
	return client
}

// apply - iterates over all opts and runs the functions to apply additional request headers
func (c *apiClient) apply(r *Request, opts []CallOption) (err error) {
	if err = r.DefaultHeaders(c.config.MasterKey); err != nil {
		return err
	}

	for i := 0; i < len(opts); i++ {
		// check to make sure someone did not pass nil ass a call option
		if opts[i] != nil {
			if err = opts[i](r); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetURI - returns a clients URI
func (c *apiClient) getURI() string {
	return c.uri
}

// GetConfig - return a clients URI
func (c *apiClient) getConfig() Config {
	return c.config
}

// EnableDebug - enables the CosmosDB debug mode
func (c *apiClient) enableDebug() {
	c.config.Debug = true
}

// DisableDebug - disables the CosmosDB debug mode
func (c *apiClient) disableDebug() {
	c.config.Debug = false
}

// Read - reads a resource by self link
func (c *apiClient) read(link string, ret interface{}, opts ...CallOption) (*Response, error) {
	return c.method("GET", link, http.StatusOK, ret, &bytes.Buffer{}, opts...)
}

// Delete - deletes a resource by self link
func (c *apiClient) delete(link string, opts ...CallOption) (*Response, error) {
	return c.method("DELETE", link, http.StatusNoContent, nil, &bytes.Buffer{}, opts...)
}

// Query - queries a resource
func (c *apiClient) query(link, query string, ret interface{}, opts ...CallOption) (*Response, error) {
	query = escapeSQL(query)
	buf := bytes.NewBufferString(querify(query))
	req, err := http.NewRequest("POST", path(c.uri, link), buf)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)
	if c.config.PartitionKeyStructField != "" {
		opts = append(opts, CrossPartition())
	}
	if err = c.apply(r, opts); err != nil {
		return nil, err
	}
	r.QueryHeaders(buf.Len())
	// revert version if collection is not partitioned
	if c.config.PartitionKeyStructField == "" {
		r.Header.Set(HeaderVersion, SupportedAPIVersionNoPartition)
	}
	// try the request and return if successful
	return c.do(r, http.StatusOK, ret)
}

// QueryWithParameters - queries a resource
func (c *apiClient) queryWithParameters(link string, query *QueryWithParameters, ret interface{}, opts ...CallOption) (*Response, error) {
	query.Query = escapeSQL(query.Query)
	q, err := stringify(query)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(q)
	req, err := http.NewRequest("POST", path(c.uri, link), buf)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)
	if c.config.PartitionKeyStructField != "" {
		opts = append(opts, CrossPartition())
	}
	if err = c.apply(r, opts); err != nil {
		return nil, err
	}
	r.QueryHeaders(buf.Len())
	// revert version if collection is not partitioned
	if c.config.PartitionKeyStructField == "" {
		r.Header.Set(HeaderVersion, SupportedAPIVersionNoPartition)
	}
	return c.do(r, http.StatusOK, ret)
}

// Create - creates a resource
func (c *apiClient) create(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("POST", link, http.StatusCreated, ret, buf, opts...)
}

// Replace - replaces a resource
func (c *apiClient) replace(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	if c.config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(body).Elem().FieldByName(c.config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	return c.method("PUT", link, http.StatusOK, ret, buf, opts...)
}

// Upsert - upserts a resource
func (c *apiClient) upsert(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	opts = append(opts, Upsert())
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	if c.config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(body).Elem().FieldByName(c.config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	return c.method(http.MethodPost, link, http.StatusOK, ret, buf, opts...)
}

// ReplaceAsync - replaces a resource
func (c *apiClient) replaceAsync(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	var Etag string
	var resource map[string]interface{}
	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, err
	}
	if valInterface, ok := resource["_etag"]; ok {
		if val, ok := valInterface.(string); ok {
			Etag = val
		}
	} else {
		return nil, errors.New("_etag does not exist for async replace")
	}
	if c.config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(body).Elem().FieldByName(c.config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	opts = append(opts, IfMatch(Etag))
	return c.method("PUT", link, http.StatusOK, ret, buf, opts...)
}

// Execute - executes a resource
func (c *apiClient) execute(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("POST", link, http.StatusOK, ret, buf, opts...)
}

// method - generic method for a resource
func (c *apiClient) method(method, link string, status int, ret interface{}, body *bytes.Buffer, opts ...CallOption) (*Response, error) {
	req, err := http.NewRequest(method, path(c.uri, link), body)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)
	if err = c.apply(r, opts); err != nil {
		return nil, err
	}
	// revert version if collection is not partitioned
	if c.config.PartitionKeyStructField == "" {
		r.Header.Set(HeaderVersion, SupportedAPIVersionNoPartition)
	}
	return c.do(r, status, ret)
}

// do - private do function
func (c *apiClient) do(r *Request, status int, data interface{}) (*Response, error) {
	if c.config.Debug && c.logger != nil {
		r.QueryMetricsHeaders()
		c.logger.Infof("CosmosDB Request: ID: %+v, Type: %+v, HTTP Request: %+v", r.rId, r.rType, r.Request)
		curl, _ := http2curl.GetCurlCommand(r.Request)
		c.logger.Infof("CURL: %s", curl)
	}
	var rr *retryablehttp.Request
	var err error
	if r.rContext != nil {
		req := r.WithContext(r.rContext)
		rr, err = retryablehttp.FromRequest(req)
	} else {
		rr, err = retryablehttp.FromRequest(r.Request)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating retryable request: %s", err)
	}
	resp, err := c.httpClient.Do(rr)
	if err != nil {
		return nil, err
	}
	if c.config.Debug && c.config.Verbose && c.logger != nil {
		c.logger.Infof("CosmosDB Request: %s", spew.Sdump(resp.Request))
		c.logger.Infof("CosmosDB Response Headers: %s", spew.Sdump(resp.Header))
		c.logger.Infof("CosmosDB Response Content-Length: %s", spew.Sdump(resp.ContentLength))
	}
	defer resp.Body.Close()
	if resp.StatusCode != status {
		err := &RequestError{}
		readJson(resp.Body, &err)
		err.StatusCode = resp.StatusCode
		err.RId = r.rId
		err.RType = r.rType
		err.Request = r.Request
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	if c.config.Debug && c.config.Verbose && c.logger != nil {
		c.logger.Infof("CosmosDB Request: %s", spew.Sdump(resp.Request))
		c.logger.Infof("CosmosDB Response Headers: %s", spew.Sdump(resp.Header))
		c.logger.Infof("CosmosDB Response Content-Length: %s", spew.Sdump(resp.ContentLength))
		c.logger.Infof("CosmosDB Response Content: %s", spew.Sdump(data))
	}
	return &Response{resp.Header}, readJson(resp.Body, data)
}

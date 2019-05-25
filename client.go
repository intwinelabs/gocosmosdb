package gocosmosdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/intwinelabs/logger"
	"github.com/moul/http2curl"
)

type Clienter interface {
	Read(link string, ret interface{}, opts ...CallOption) (*Response, error)
	Delete(link string, opts ...CallOption) (*Response, error)
	Query(link string, query string, ret interface{}, opts ...CallOption) (*Response, error)
	QueryWithParameters(link string, query *QueryWithParameters, ret interface{}, opts ...CallOption) (*Response, error)
	Create(link string, body, ret interface{}, opts ...CallOption) (*Response, error)
	Upsert(link string, body, ret interface{}, opts ...CallOption) (*Response, error)
	Replace(link string, body, ret interface{}, opts ...CallOption) (*Response, error)
	ReplaceAsync(link string, body, ret interface{}, opts ...CallOption) (*Response, error)
	Execute(link string, body, ret interface{}, opts ...CallOption) (*Response, error)
	GetURI() string
	GetConfig() Config
	EnableDebug()
	DisableDebug()
}

type Client struct {
	Url    string
	Config Config
	http.Client
	Logger *logger.Logger
}

func (c *Client) apply(r *Request, opts []CallOption) (err error) {
	if err = r.DefaultHeaders(c.Config.MasterKey); err != nil {
		return err
	}

	for i := 0; i < len(opts); i++ {
		if err = opts[i](r); err != nil {
			return err
		}
	}
	return nil
}

// GetURI return a clients URI
func (c *Client) GetURI() string {
	return c.Url
}

// GetConfig return a clients URI
func (c *Client) GetConfig() Config {
	return c.Config
}

// EnableDebug enables the CosmosDB debug in config
func (c *Client) EnableDebug() {
	c.Config.Debug = true
}

// DisableDebug disables the CosmosDB debug in config
func (c *Client) DisableDebug() {
	c.Config.Debug = false
}

// Read resource by self link
func (c *Client) Read(link string, ret interface{}, opts ...CallOption) (*Response, error) {
	return c.method("GET", link, http.StatusOK, ret, &bytes.Buffer{}, opts...)
}

// Delete resource by self link
func (c *Client) Delete(link string, opts ...CallOption) (*Response, error) {
	return c.method("DELETE", link, http.StatusNoContent, nil, &bytes.Buffer{}, opts...)
}

// Query resource
func (c *Client) Query(link, query string, ret interface{}, opts ...CallOption) (*Response, error) {
	query = escapeSQL(query)
	buf := bytes.NewBufferString(querify(query))
	req, err := http.NewRequest("POST", path(c.Url, link), buf)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)
	if c.Config.PartitionKeyStructField != "" {
		opts = append(opts, CrossPartition())
	}
	if err = c.apply(r, opts); err != nil {
		return nil, err
	}
	r.QueryHeaders(buf.Len())
	return c.do(r, http.StatusOK, ret)
}

// QueryWithParameters resource
func (c *Client) QueryWithParameters(link string, query *QueryWithParameters, ret interface{}, opts ...CallOption) (*Response, error) {
	query.Query = escapeSQL(query.Query)
	q, err := stringify(query)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(q)
	req, err := http.NewRequest("POST", path(c.Url, link), buf)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)
	if c.Config.PartitionKeyStructField != "" {
		opts = append(opts, CrossPartition())
	}
	if err = c.apply(r, opts); err != nil {
		return nil, err
	}
	r.QueryHeaders(buf.Len())
	return c.do(r, http.StatusOK, ret)
}

// Create resource
func (c *Client) Create(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("POST", link, http.StatusCreated, ret, buf, opts...)
}

// Replace resource
func (c *Client) Replace(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	if c.Config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(body).Elem().FieldByName(c.Config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	return c.method("PUT", link, http.StatusOK, ret, buf, opts...)
}

// Upsert resource
func (c *Client) Upsert(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	opts = append(opts, Upsert())
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	if c.Config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(body).Elem().FieldByName(c.Config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	return c.method(http.MethodPost, link, http.StatusOK, ret, buf, opts...)
}

// ReplaceAsync resource
func (c *Client) ReplaceAsync(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
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
	if c.Config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(body).Elem().FieldByName(c.Config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	opts = append(opts, IfMatch(Etag))
	return c.method("PUT", link, http.StatusOK, ret, buf, opts...)
}

// Replace resource
func (c *Client) Execute(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("POST", link, http.StatusOK, ret, buf, opts...)
}

// Private generic method resource
func (c *Client) method(method, link string, status int, ret interface{}, body *bytes.Buffer, opts ...CallOption) (*Response, error) {
	req, err := http.NewRequest(method, path(c.Url, link), body)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)
	if err = c.apply(r, opts); err != nil {
		return nil, err
	}
	return c.do(r, status, ret)
}

// Private Do function, DRY
func (c *Client) do(r *Request, status int, data interface{}) (*Response, error) {
	if c.Config.Debug {
		r.QueryMetricsHeaders()
		c.Logger.Infof("CosmosDB Request: ID: %+v, Type: %+v, HTTP Request: %+v", r.rId, r.rType, r.Request)
		curl, _ := http2curl.GetCurlCommand(r.Request)
		c.Logger.Infof("CURL: %s", curl)
	}
	resp, err := c.Do(r.Request)
	if err != nil {
		return nil, fmt.Errorf("Request: Id: %+v, Type: %+v, HTTP: %+v, Error: %s", r.rId, r.rType, r.Request, err)
	}
	if c.Config.Debug && c.Config.Verbose {
		c.Logger.Infof("CosmosDB Request: %s", spew.Sdump(resp.Request))
		c.Logger.Infof("CosmosDB Response Headers: %s", spew.Sdump(resp.Header))
		c.Logger.Infof("CosmosDB Response Content-Length: %s", spew.Sdump(resp.Header))
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
	if c.Config.Debug && c.Config.Verbose {
		c.Logger.Infof("CosmosDB Request: %s", spew.Sdump(resp.Request))
		c.Logger.Infof("CosmosDB Response Headers: %s", spew.Sdump(resp.Header))
		c.Logger.Infof("CosmosDB Response Content-Length: %s", spew.Sdump(resp.Header))
		c.Logger.Infof("CosmosDB Response Content: %s", spew.Sdump(data))
	}
	return &Response{resp.Header}, readJson(resp.Body, data)
}

// Generate link
func path(url string, args ...string) (link string) {
	args = append([]string{url}, args...)
	link = strings.Join(args, "/")
	return
}

// Read json response to given interface(struct, map, ..)
func readJson(reader io.Reader, data interface{}) error {
	return json.NewDecoder(reader).Decode(&data)
}

// Stringify query-string as CosmosDB expected
func querify(query string) string {
	return fmt.Sprintf(`{ "%s": "%s" }`, "query", query)
}

// Stringify body data
func stringify(body interface{}) (bt []byte, err error) {
	switch t := body.(type) {
	case string:
		bt = []byte(t)
	case []byte:
		bt = t
	default:
		bt, err = json.Marshal(t)
	}
	return
}

package gocosmosdb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/intwinelabs/logger"
	"github.com/stretchr/testify/assert"
)

var log = logger.New()
var httpClient = retryablehttp.NewClient()

func init() {
	httpClient.RetryMax = 0
}

type RequestRecorder struct {
	Header http.Header
	Body   string
}

type MockServer struct {
	*httptest.Server
	RequestRecorder
	Status interface{}
}

func (m *MockServer) SetStatus(status int) {
	m.Status = status
}

func (s *MockServer) Record(r *http.Request) {
	s.Header = r.Header
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	s.Body = string(b)
}

func (s *MockServer) AssertHeaders(t *testing.T, headers ...string) {
	assert := assert.New(t)
	for _, k := range headers {
		assert.NotNil(s.Header[k])
	}
}

func ServerFactory(resp ...interface{}) *MockServer {
	s := &MockServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		// Record the last request
		s.Record(r)
		if v, ok := resp[0].(int); ok {
			err := fmt.Errorf(`{"code": "500", "message": "CosmosDB error"}`)
			http.Error(w, err.Error(), v)
		} else {
			if status, ok := s.Status.(int); ok {
				w.WriteHeader(status)
			}
			fmt.Fprintln(w, resp[0])
		}
		resp = resp[1:]
	}))
	return s
}

func TestGetURI(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		logger: log,
	}

	// First call
	uri := client.getURI()
	assert.Equal(s.URL, uri)
}

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		logger: log,
	}

	// First call
	expConf := Config{
		MasterKey: "YXJpZWwNCg==",
	}
	conf := client.getConfig()
	assert.Equal(expConf, conf)
}

func TestRead(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		httpClient: httpClient,
		logger:     log,
	}

	// First call
	var db Database
	_, err := client.read("dbs/b7NTAS==/", &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second Call, when StatusCode != StatusOK
	_, err = client.read("dbs/b7NCAA==/colls/Ad352/", &db)
	assert.Contains(err.Error(), "giving up after 1 attempts")
}

func TestQuery(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		httpClient: httpClient,
		logger:     log,
	}

	// First call
	var db Database
	_, err := client.query("dbs", "SELECT * FROM ROOT r", &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	s.AssertHeaders(t, HeaderContentLength, HeaderContentType, HeaderIsQuery)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second Call, when StatusCode != StatusOK
	_, err = client.read("/dbs/b7NCAA==/colls/Ad352/", &db)
	assert.Contains(err.Error(), "giving up after 1 attempts")
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, `{"id": "9"}`, 500)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		httpClient: httpClient,
		logger:     log,
	}

	// First call
	var db Database
	_, err := client.create("dbs", `{"id": 3}`, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second call
	var doc, tDoc Document
	tDoc.Id = "9"
	_, err = client.create("dbs", tDoc, &doc)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(doc.Id, tDoc.Id, "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Last Call, when StatusCode != StatusOK && StatusCreated
	_, err = client.create("dbs", tDoc, &doc)
	assert.Contains(err.Error(), "giving up after 1 attempts")
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`10`, 500)
	s.SetStatus(http.StatusNoContent)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		httpClient: httpClient,
		logger:     log,
	}

	// First call
	_, err := client.delete("dbs/b7NTAS==/")
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Nil(err, "err should be nil")

	// Second Call, when StatusCode != StatusOK
	_, err = client.delete("dbs/b7NCAA==/colls/Ad352/")
	assert.Contains(err.Error(), "giving up after 1 attempts")
}

func TestReplace(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, `{"id": "9"}`, 500)
	s.SetStatus(http.StatusOK)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		httpClient: httpClient,
		logger:     log,
	}

	// First call
	var db Database
	_, err := client.replace("dbs", `{"id": 3}`, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second call
	var doc, tDoc Document
	tDoc.Id = "9"
	_, err = client.replace("dbs", tDoc, &doc)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(doc.Id, tDoc.Id, "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Last Call, when StatusCode != StatusOK && StatusCreated
	_, err = client.replace("dbs", tDoc, &doc)
	assert.Contains(err.Error(), "giving up after 1 attempts")
}

func TestExecute(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, `{"id": "9"}`, 500)
	s.SetStatus(http.StatusOK)
	defer s.Close()
	client := &apiClient{
		uri: s.URL,
		config: Config{
			MasterKey: "YXJpZWwNCg==",
		},
		httpClient: httpClient,
		logger:     log,
	}

	// First call
	var db Database
	_, err := client.execute("dbs", `{"id": 3}`, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second call
	var doc, tDoc Document
	tDoc.Id = "9"
	_, err = client.execute("dbs", tDoc, &doc)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(doc.Id, tDoc.Id, "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Last Call, when StatusCode != StatusOK && StatusCreated
	_, err = client.execute("dbs", tDoc, &doc)
	assert.Contains(err.Error(), "giving up after 1 attempts")
}

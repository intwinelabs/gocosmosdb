//
// This library started as a fork of `github.com/nerdylikeme/go-documentdb`
//

package gocosmosdb

import (
	"errors"
	"reflect"
	"time"

	"github.com/intwinelabs/logger"
)

// Config - Stores configuration for the gocosmosdb client
type Config struct {
	MasterKey               string
	Debug                   bool
	Verbose                 bool
	PartitionKeyStructField string // eg. "Id"
	PartitionKeyPath        string // slash denoted path eg. "/id"
	RetryWaitMin            time.Duration
	RetryWaitMax            time.Duration
	RetryMax                int
	Pooled                  bool
}

// CosmosDB - Struct that stores the client and logger
type CosmosDB struct {
	client *Client
	Config Config
	Logger *logger.Logger
}

// New - Creates CosmosDB Client and returns it
func New(url string, config Config, log *logger.Logger) *CosmosDB {
	client := &Client{}
	client.URI = url
	client.Config = config
	client.Logger = log
	return &CosmosDB{client, config, log}
}

// GetURI - returns the CosmosDB URI
func (c *CosmosDB) GetURI() string {
	return c.client.GetURI()
}

// GetConfig - returns the CosmosDB config
func (c *CosmosDB) GetConfig() Config {
	return c.client.GetConfig()
}

// EnableDebug - enables the CosmosDB debugging
func (c *CosmosDB) EnableDebug() {
	c.client.EnableDebug()
}

// DisableDebug - disables the CosmosDB debugging
func (c *CosmosDB) DisableDebug() {
	c.client.DisableDebug()
}

// ReadDatabase - Retrieves a database resource by performing a GET on the database resource.
//	db, err := client.ReadDatabase("https://{databaseaccount}.documents.azure.com/dbs/{db-id}")
func (c *CosmosDB) ReadDatabase(link string, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Read(link, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadCollection - Retrieves a collection by performing a GET on a specific collection resource.
//	coll, err := client.ReadCollection("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/colls/{coll-id}")
func (c *CosmosDB) ReadCollection(link string, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.Read(link, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadDocument - Retrieves a document by performing a GET on a specific document resource and marshales the document into that passed docStruct
//	err = client.ReadDocument("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/colls/{coll-id}/docs/{doc-id}", &docStruct)
func (c *CosmosDB) ReadDocument(link string, doc interface{}, opts ...CallOption) (err error) {
	_, err = c.client.Read(link, &doc, opts...)
	return
}

// ReadStoredProcedure - Retrieves a stored procedure by performing a GET on a specific stored procedure resource.
// sproc, err := client.ReadStoredProcedure("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/sprocs/{sproc-id}")
func (c *CosmosDB) ReadStoredProcedure(link string, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Read(link, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadUserDefinedFunction - Retrieves a user defined function by performing a GET on a specific user defined function resource.
// udf, err := client.ReadUserDefinedFunction("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/udfs/{udf-id}")
func (c *CosmosDB) ReadUserDefinedFunction(link string, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Read(link, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadDatabases - Retrieves all databases by performing a GET on a specific account.
//	dbs, err := client.ReadDatabases("https://{databaseaccount}.documents.azure.com/dbs")
func (c *CosmosDB) ReadDatabases(opts ...CallOption) (dbs []Database, err error) {
	return c.QueryDatabases("", opts...)
}

// ReadCollections - Retrieves all collections by performing a GET on a specific database.
//	colls, err := client.ReadCollections("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/colls")
func (c *CosmosDB) ReadCollections(db string, opts ...CallOption) (colls []Collection, err error) {
	return c.QueryCollections(db, "", opts...)
}

// ReadStoredProcedures - Retrieves all stored procedures by performing a GET on a specific database.
//	sprocs, err := client.ReadStoredProcedures("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/sprocs")
func (c *CosmosDB) ReadStoredProcedures(coll string, opts ...CallOption) (sprocs []Sproc, err error) {
	return c.QueryStoredProcedures(coll, "", opts...)
}

// ReadUserDefinedFunctions - Retrieves all user defined functions by performing a GET on a specific database.
//	udfs, err := client.ReadUserDefinedFunctions("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/udfs")
func (c *CosmosDB) ReadUserDefinedFunctions(coll string, opts ...CallOption) (udfs []UDF, err error) {
	return c.QueryUserDefinedFunctions(coll, "", opts...)
}

// ReadDocuments - Retrieves a stored procedure by performing a GET on a specific stored procedure resource.
//	err = client.ReadDocuments("https://{databaseaccount}.documents.azure.com/dbs/{db-id}/colls/{coll-id}/docs", &docStructSlice)
func (c *CosmosDB) ReadDocuments(coll string, docs interface{}, opts ...CallOption) (err error) {
	return c.QueryDocuments(coll, "", docs, opts...)
}

// QueryDatabases - Retrieves all databases that satisfy the passed query.
//	dbs, err := client.QueryDatabases("SELECT * FROM ROOT r")
func (c *CosmosDB) QueryDatabases(query string, opts ...CallOption) (dbs []Database, err error) {
	data := struct {
		Databases []Database `json:"Databases,omitempty"`
		Count     int        `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		_, err = c.client.Query("dbs", query, &data, opts...)
	} else {
		_, err = c.client.Read("dbs", &data, opts...)
	}
	if dbs = data.Databases; err != nil {
		dbs = nil
	}
	return
}

// QueryCollections - Retrieves all collections that satisfy passed query.
//	colls, err := client.QueryCollections("SELECT * FROM ROOT r")
func (c *CosmosDB) QueryCollections(db, query string, opts ...CallOption) (colls []Collection, err error) {
	data := struct {
		Collections []Collection `json:"DocumentCollections,omitempty"`
		Count       int          `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		_, err = c.client.Query(db+"colls/", query, &data, opts...)
	} else {
		_, err = c.client.Read(db+"colls/", &data, opts...)
	}
	if colls = data.Collections; err != nil {
		colls = nil
	}
	return
}

// QueryStoredProcedures - Retrieves all stored procedures that satisfy the passed query.
//	colls, err := client.QueryStoredProcedures("SELECT * FROM ROOT r")
func (c *CosmosDB) QueryStoredProcedures(coll, query string, opts ...CallOption) (sprocs []Sproc, err error) {
	data := struct {
		Sprocs []Sproc `json:"StoredProcedures,omitempty"`
		Count  int     `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		_, err = c.client.Query(coll+"sprocs/", query, &data, opts...)
	} else {
		_, err = c.client.Read(coll+"sprocs/", &data, opts...)
	}
	if sprocs = data.Sprocs; err != nil {
		sprocs = nil
	}
	return
}

// QueryUserDefinedFunctions - Retrieves all user defined functions that satisfy the passed query.
//	colls, err := client.QueryUserDefinedFunctions("SELECT * FROM ROOT r")
func (c *CosmosDB) QueryUserDefinedFunctions(coll, query string, opts ...CallOption) (udfs []UDF, err error) {
	data := struct {
		Udfs  []UDF `json:"UserDefinedFunctions,omitempty"`
		Count int   `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		_, err = c.client.Query(coll+"udfs/", query, &data, opts...)
	} else {
		_, err = c.client.Read(coll+"udfs/", &data, opts...)
	}
	if udfs = data.Udfs; err != nil {
		udfs = nil
	}
	return
}

// QueryDocuments - Retrieves all documents in a collection that satisfy the passed query.
//	err := client.QueryDocuments(coll, "SELECT * FROM ROOT r", &docs)
func (c *CosmosDB) QueryDocuments(coll, query string, docs interface{}, opts ...CallOption) (err error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if len(query) > 0 {
		_, err = c.client.Query(coll+"docs/", query, &data, opts...)
	} else {
		_, err = c.client.Read(coll+"docs/", &data)
	}
	return
}

// QueryDocumentsWithParameters - Retrieves all documents in a collection that satisfy a passed query with parameters.
//	colls, err := client.QueryCollections(coll, queryWithParams, &docs)
func (c *CosmosDB) QueryDocumentsWithParameters(coll string, query *QueryWithParameters, docs interface{}, opts ...CallOption) (err error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if query != nil {
		_, err = c.client.QueryWithParameters(coll+"docs/", query, &data, opts...)
	} else {
		err = errors.New("QueryWithParameters cannot be nil")
	}
	return
}

// QueryPartitionKeyRanges - Retrieves all partition ranges in a collection.
//	pks, err := client.QueryCollections(coll, "SELECT * FROM ROOT r")
func (c *CosmosDB) QueryPartitionKeyRanges(coll string, query string, opts ...CallOption) (ranges []PartitionKeyRange, err error) {
	data := queryPartitionKeyRangesRequest{}
	if len(query) > 0 {
		_, err = c.client.Query(coll+"pkranges/", query, &data, opts...)
	} else {
		_, err = c.client.Read(coll+"pkranges/", &data, opts...)
	}
	if ranges = data.Ranges; err != nil {
		ranges = nil
	}
	return
}

// CreateDatabase -
func (c *CosmosDB) CreateDatabase(body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Create("dbs", body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateCollection -
func (c *CosmosDB) CreateCollection(db string, body interface{}, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.Create(db+"colls/", body, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateStoredProcedure -
func (c *CosmosDB) CreateStoredProcedure(coll string, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Create(coll+"sprocs/", body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateUserDefinedFunction -
func (c *CosmosDB) CreateUserDefinedFunction(coll string, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Create(coll+"udfs/", body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateDocument -
func (c *CosmosDB) CreateDocument(coll string, doc interface{}, opts ...CallOption) (*Response, error) {
	id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	if id.IsValid() && id.CanSet() && id.String() == "" {
		id.SetString(genId())
	}
	if c.Config.PartitionKeyStructField != "" {
		partKey := reflect.ValueOf(doc).Elem().FieldByName(c.Config.PartitionKeyStructField)
		partKeyI := partKey.Interface()
		opts = append(opts, PartitionKey(partKeyI))
	}
	return c.client.Create(coll+"docs/", doc, &doc, opts...)
}

// UpsertDocument -
func (c *CosmosDB) UpsertDocument(coll string, doc interface{}, opts ...CallOption) (*Response, error) {
	id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	if id.IsValid() && id.CanSet() && id.String() == "" {
		id.SetString(genId())
	}
	return c.client.Upsert(coll+"docs/", doc, &doc, opts...)
}

// DeleteDatabase -
func (c *CosmosDB) DeleteDatabase(link string) (*Response, error) {
	return c.client.Delete(link)
}

// DeleteCollection -
func (c *CosmosDB) DeleteCollection(link string) (*Response, error) {
	return c.client.Delete(link)
}

// DeleteDocument -
func (c *CosmosDB) DeleteDocument(link string, opts ...CallOption) (*Response, error) {
	return c.client.Delete(link, opts...)
}

// DeleteStoredProcedure -
func (c *CosmosDB) DeleteStoredProcedure(link string) (*Response, error) {
	return c.client.Delete(link)
}

// DeleteUserDefinedFunction -
func (c *CosmosDB) DeleteUserDefinedFunction(link string) (*Response, error) {
	return c.client.Delete(link)
}

// ReplaceDatabase -
func (c *CosmosDB) ReplaceDatabase(link string, body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Replace(link, body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReplaceDocument -
func (c *CosmosDB) ReplaceDocument(link string, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.Replace(link, doc, &doc, opts...)
}

// ReplaceDocumentAsync -
func (c *CosmosDB) ReplaceDocumentAsync(link string, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.ReplaceAsync(link, doc, &doc, opts...)
}

// ReplaceStoredProcedure -
func (c *CosmosDB) ReplaceStoredProcedure(link string, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Replace(link, body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReplaceUserDefinedFunction -
func (c *CosmosDB) ReplaceUserDefinedFunction(link string, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Replace(link, body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ExecuteStoredProcedure -
func (c *CosmosDB) ExecuteStoredProcedure(link string, params, body interface{}, opts ...CallOption) (err error) {
	_, err = c.client.Execute(link, params, &body, opts...)
	return
}

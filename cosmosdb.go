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
	client *apiClient
	Config Config
	Logger *logger.Logger
}

// New - Creates CosmosDB Client and returns it
func New(url string, config Config, log *logger.Logger) *CosmosDB {
	client := newAPIClient(&config)
	client.uri = url
	client.config = config
	client.logger = log
	return &CosmosDB{client, config, log}
}

// GetURI - returns the CosmosDB URI
func (c *CosmosDB) GetURI() string {
	return c.client.getURI()
}

// GetConfig - returns the CosmosDB config
func (c *CosmosDB) GetConfig() Config {
	return c.client.getConfig()
}

// EnableDebug - enables the CosmosDB debugging
func (c *CosmosDB) EnableDebug() {
	c.client.enableDebug()
}

// DisableDebug - disables the CosmosDB debugging
func (c *CosmosDB) DisableDebug() {
	c.client.disableDebug()
}

// ReadDatabase - Retrieves a database resource by performing a GET on the database resource.
//	db, err := client.ReadDatabase("dbs/{db-id}")
func (c *CosmosDB) ReadDatabase(link string, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.read(link, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadCollection - Retrieves a collection by performing a GET on a specific collection resource.
//	coll, err := client.ReadCollection("dbs/{db-id}/colls/{coll-id}")
func (c *CosmosDB) ReadCollection(link string, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.read(link, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadDocument - Retrieves a document by performing a GET on a specific document resource and marshals the document into that passed docStruct
//	err = client.ReadDocument("dbs/{db-id}/colls/{coll-id}/docs/{doc-id}", &docStruct)
func (c *CosmosDB) ReadDocument(link string, doc interface{}, opts ...CallOption) (err error) {
	_, err = c.client.read(link, &doc, opts...)
	return
}

// ReadStoredProcedure - Retrieves a stored procedure by performing a GET on a specific stored procedure resource.
// sproc, err := client.ReadStoredProcedure("dbs/{db-id}/sprocs/{sproc-id}")
func (c *CosmosDB) ReadStoredProcedure(link string, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.read(link, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadUserDefinedFunction - Retrieves a user defined function by performing a GET on a specific user defined function resource.
// udf, err := client.ReadUserDefinedFunction("dbs/{db-id}/udfs/{udf-id}")
func (c *CosmosDB) ReadUserDefinedFunction(link string, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.read(link, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReadDatabases - Retrieves all databases by performing a GET on a specific account.
//	dbs, err := client.ReadDatabases("dbs")
func (c *CosmosDB) ReadDatabases(opts ...CallOption) (dbs []Database, err error) {
	return c.QueryDatabases("", opts...)
}

// ReadCollections - Retrieves all collections by performing a GET on a specific database.
//	colls, err := client.ReadCollections("dbs/{db-id}/colls")
func (c *CosmosDB) ReadCollections(db string, opts ...CallOption) (colls []Collection, err error) {
	return c.QueryCollections(db, "", opts...)
}

// ReadStoredProcedures - Retrieves all stored procedures by performing a GET on a specific database.
//	sprocs, err := client.ReadStoredProcedures("dbs/{db-id}/sprocs")
func (c *CosmosDB) ReadStoredProcedures(coll string, opts ...CallOption) (sprocs []Sproc, err error) {
	return c.QueryStoredProcedures(coll, "", opts...)
}

// ReadUserDefinedFunctions - Retrieves all user defined functions by performing a GET on a specific database.
//	udfs, err := client.ReadUserDefinedFunctions("dbs/{db-id}/udfs")
func (c *CosmosDB) ReadUserDefinedFunctions(coll string, opts ...CallOption) (udfs []UDF, err error) {
	return c.QueryUserDefinedFunctions(coll, "", opts...)
}

// ReadDocuments - Retrieves a stored procedure by performing a GET on a specific stored procedure resource.
//	err = client.ReadDocuments("dbs/{db-id}/colls/{coll-id}/docs", &docStructSlice)
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
		_, err = c.client.query("dbs", query, &data, opts...)
	} else {
		_, err = c.client.read("dbs", &data, opts...)
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
		_, err = c.client.query(db+"colls/", query, &data, opts...)
	} else {
		_, err = c.client.read(db+"colls/", &data, opts...)
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
		_, err = c.client.query(coll+"sprocs/", query, &data, opts...)
	} else {
		_, err = c.client.read(coll+"sprocs/", &data, opts...)
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
		_, err = c.client.query(coll+"udfs/", query, &data, opts...)
	} else {
		_, err = c.client.read(coll+"udfs/", &data, opts...)
	}
	if udfs = data.Udfs; err != nil {
		udfs = nil
	}
	return
}

// QueryDocuments - Retrieves all documents in a collection that satisfy the passed query and marshals them into the passed interface.
//	err := client.QueryDocuments(coll, "SELECT * FROM ROOT r", &docs)
func (c *CosmosDB) QueryDocuments(coll, query string, docs interface{}, opts ...CallOption) (err error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if len(query) > 0 {
		_, err = c.client.query(coll+"docs/", query, &data, opts...)
	} else {
		_, err = c.client.read(coll+"docs/", &data)
	}
	return
}

// QueryDocumentsWithParameters - Retrieves all documents in a collection that satisfy a passed query with parameters and marshals them into the passed interface.
//	err := client.QueryDocumentsWithParameters(coll, queryWithParams, &docs)
func (c *CosmosDB) QueryDocumentsWithParameters(coll string, query *QueryWithParameters, docs interface{}, opts ...CallOption) (err error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if query != nil {
		_, err = c.client.queryWithParameters(coll+"docs/", query, &data, opts...)
	} else {
		err = errors.New("QueryWithParameters cannot be nil")
	}
	return
}

// QueryPartitionKeyRanges - Retrieves all partition ranges in a collection.
//	pks, err := client.QueryPartitionKeyRanges(coll, "SELECT * FROM ROOT r")
func (c *CosmosDB) QueryPartitionKeyRanges(coll string, query string, opts ...CallOption) (ranges []PartitionKeyRange, err error) {
	data := struct {
		PartitionKeyRanges []PartitionKeyRange `json:"PartitionKeyRanges,omitempty"`
		Count              int                 `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		_, err = c.client.query(coll+"pkranges/", query, &data, opts...)
	} else {
		_, err = c.client.read(coll+"pkranges/", &data, opts...)
	}
	if ranges = data.PartitionKeyRanges; err != nil {
		ranges = nil
	}
	return
}

// CreateDatabase - Creates a new database in the database account.
//	db, err := client.CreateDatabase(`{ "id": "db-id" }`)
func (c *CosmosDB) CreateDatabase(body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.create("dbs", body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateCollection - Creates a new collections in the database.
//	coll, err := client.CreateCollection("dbs/{db-id}/", `{"id": "coll-id"}`)
func (c *CosmosDB) CreateCollection(db string, body interface{}, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.create(db+"colls/", body, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateStoredProcedure - Creates a new stored procedure in the collection.
//	sprocBody := gocosmosdb.Sproc{
//    	Body: "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World\");\r\n}",
//    	Id: "sproc_1"
//	}
//	sproc, err := client.CreateStoredProcedure("dbs/{db-id}/colls/{coll-id}/sprocs", &sprocBody)
func (c *CosmosDB) CreateStoredProcedure(coll string, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.create(coll+"sprocs/", body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateUserDefinedFunction - Creates a new user defined function in the collection.
//	udfBody := gocosmosdb.UDF{
//    	Body: "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",
//    	Id: "simpleTaxUDF"
//	}
//	udf, err := client.CreateUserDefinedFunction("dbs/{db-id}/colls/{coll-id}/udfs", &udfBody)
func (c *CosmosDB) CreateUserDefinedFunction(coll string, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.create(coll+"udfs/", body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// CreateDocument - Creates a new document in the collection.
//	err := client.CreateDocument("dbs/{db-id}/colls/{coll-id}", &doc)
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
	return c.client.create(coll+"docs/", doc, &doc, opts...)
}

// UpsertDocument - Creates a new document or replaces the existing document with matching id in the collection.
//	err := client.UpsertDocument("dbs/{db-id}/colls/{coll-id}", &doc)
func (c *CosmosDB) UpsertDocument(coll string, doc interface{}, opts ...CallOption) (*Response, error) {
	id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	if id.IsValid() && id.CanSet() && id.String() == "" {
		id.SetString(genId())
	}
	return c.client.upsert(coll+"docs/", doc, &doc, opts...)
}

// DeleteDatabase - Deletes a database from a database account.
//	err := client.DeleteDatabase("dbs/{db-id}")
func (c *CosmosDB) DeleteDatabase(link string) (*Response, error) {
	return c.client.delete(link)
}

// DeleteCollection - Deletes a collection from a database.
//	err := client.DeleteCollection("dbs/{db-id}/colls/{coll-id}")
func (c *CosmosDB) DeleteCollection(link string) (*Response, error) {
	return c.client.delete(link)
}

// DeleteDocument -  Deletes a document from a collection.
//	err := client.DeleteDocument("dbs/{db-id}/colls/{coll-id}/docs/{doc-id}")
func (c *CosmosDB) DeleteDocument(link string, opts ...CallOption) (*Response, error) {
	return c.client.delete(link, opts...)
}

// DeleteStoredProcedure -  Deletes a stored procedure from a collection.
//	err := client.DeleteStoredProcedure("dbs/{db-id}/colls/{coll-id}/sprocs/{sproc-id}")
func (c *CosmosDB) DeleteStoredProcedure(link string) (*Response, error) {
	return c.client.delete(link)
}

// DeleteUserDefinedFunction -  Deletes a user defined function from a collection.
//	err := client.DeleteUserDefinedFunction("dbs/{db-id}/colls/{coll-id}/udfs/{udf-id}")
func (c *CosmosDB) DeleteUserDefinedFunction(link string) (*Response, error) {
	return c.client.delete(link)
}

// ReplaceDatabase - Replaces a existing database in a database account.
//	db, err := client.ReplaceDatabase("dbs/{db-id}", "`{ "id": "new-db-id" }`)
func (c *CosmosDB) ReplaceDatabase(link string, body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.replace(link, body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReplaceDocument - Replaces a existing document in a collection.
//	db, err := client.ReplaceDocument("dbs/{db-id}/colls/{coll-id}/docs/{doc-id}", &doc)
func (c *CosmosDB) ReplaceDocument(link string, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.replace(link, doc, &doc, opts...)
}

// ReplaceDocumentAsync - Replaces a document that has a matching etag.
//	db, err := client.ReplaceDocumentAsync("dbs/{db-id}/colls/{coll-id}/docs/{doc-id}", &doc)
func (c *CosmosDB) ReplaceDocumentAsync(link string, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.replaceAsync(link, doc, &doc, opts...)
}

// ReplaceStoredProcedure - Replaces a stored procedure in a collection.
//	db, err := client.ReplaceDatabase("dbs/{db-id}/colls/{coll-id}/sprocs/{sproc-id}", &sprocBody)
func (c *CosmosDB) ReplaceStoredProcedure(link string, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.replace(link, body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ReplaceUserDefinedFunction - Replaces a user defined function in a collection.
//	db, err := client.ReplaceDatabase("dbs/{db-id}/colls/{coll-id}/udfs/{udf-id}", &udfBody)
func (c *CosmosDB) ReplaceUserDefinedFunction(link string, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.replace(link, body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// ExecuteStoredProcedure - Executes a stored procedure and marshals the data into the passed interface.
//	err := client.ExecuteStoredProcedure("dbs/{db-id}/colls/{coll-id}/sprocs/{sproc-id}", []interface{}{p1, p2}, &docs)
func (c *CosmosDB) ExecuteStoredProcedure(link string, params, body interface{}, opts ...CallOption) (err error) {
	_, err = c.client.execute(link, params, &body, opts...)
	return
}

//
// This library started as a fork of `github.com/nerdylikeme/go-documentdb`
//

package gocosmosdb

import (
	"errors"
	"reflect"

	"github.com/intwinelabs/logger"
)

type Config struct {
	MasterKey               string
	Debug                   bool
	Verbose                 bool
	PartitionKeyStructField string // eg. "Id"
	PartitionKeyPath        string // slash denoted path eg. "/id"
}

type CosmosDB struct {
	client Clienter
	Logger *logger.Logger
}

// Create CosmosDBClient
func New(url string, config Config, log *logger.Logger) *CosmosDB {
	client := &Client{}
	client.Url = url
	client.Config = config
	client.Logger = log
	return &CosmosDB{client, log}
}

// GetURI returns the CosmosDB URI
func (c *CosmosDB) GetURI() string {
	return c.client.GetURI()
}

// GetConfig return the CosmosDB config
func (c *CosmosDB) GetConfig() Config {
	return c.client.GetConfig()
}

// EnableDebug enables the CosmosDB debug in config
func (c *CosmosDB) EnableDebug() {
	c.client.EnableDebug()
}

// DisableDebug disables the CosmosDB debug in config
func (c *CosmosDB) DisableDebug() {
	c.client.DisableDebug()
}

// Read database by self link
func (c *CosmosDB) ReadDatabase(link string, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Read(link, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read collection by self link
func (c *CosmosDB) ReadCollection(link string, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.Read(link, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read document by self link
func (c *CosmosDB) ReadDocument(link string, doc interface{}, opts ...CallOption) (err error) {
	_, err = c.client.Read(link, &doc, opts...)
	return
}

// Read sporc by self link
func (c *CosmosDB) ReadStoredProcedure(link string, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Read(link, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read udf by self link
func (c *CosmosDB) ReadUserDefinedFunction(link string, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Read(link, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read all databases
func (c *CosmosDB) ReadDatabases(opts ...CallOption) (dbs []Database, err error) {
	return c.QueryDatabases("", opts...)
}

// Read all collections by db selflink
func (c *CosmosDB) ReadCollections(db string, opts ...CallOption) (colls []Collection, err error) {
	return c.QueryCollections(db, "", opts...)
}

// Read all sprocs by collection self link
func (c *CosmosDB) ReadStoredProcedures(coll string, opts ...CallOption) (sprocs []Sproc, err error) {
	return c.QueryStoredProcedures(coll, "", opts...)
}

// Read all udfs by collection self link
func (c *CosmosDB) ReadUserDefinedFunctions(coll string, opts ...CallOption) (udfs []UDF, err error) {
	return c.QueryUserDefinedFunctions(coll, "", opts...)
}

// Read all collection documents by self link
// TODO: use iterator for heavy transactions
func (c *CosmosDB) ReadDocuments(coll string, docs interface{}, opts ...CallOption) (err error) {
	return c.QueryDocuments(coll, "", docs, opts...)
}

// Read all databases that satisfy a query
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

// Read all db-collection that satisfy a query
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

// Read all collection `sprocs` that satisfy a query
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

// Read all collection `udfs` that satisfy a query
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

// Read all documents in a collection that satisfy a query
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

// Read all documents in a collection that satisfy a query with parameters
func (c *CosmosDB) QueryDocumentsWithParmeters(coll string, query *QueryWithParameters, docs interface{}, opts ...CallOption) (err error) {
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

// Read collection's partition ranges
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

// Create database
func (c *CosmosDB) CreateDatabase(body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Create("dbs", body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create collection
func (c *CosmosDB) CreateCollection(db string, body interface{}, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.Create(db+"colls/", body, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create stored procedure
func (c *CosmosDB) CreateStoredProcedure(coll string, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Create(coll+"sprocs/", body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create user defined function
func (c *CosmosDB) CreateUserDefinedFunction(coll string, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Create(coll+"udfs/", body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create document
func (c *CosmosDB) CreateDocument(coll string, doc interface{}, opts ...CallOption) (*Response, error) {
	id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	if id.IsValid() && id.CanSet() && id.String() == "" {
		id.SetString(genId())
	}
	return c.client.Create(coll+"docs/", doc, &doc, opts...)
}

// Upsert document
func (c *CosmosDB) UpsertDocument(coll string, doc interface{}, opts ...CallOption) (*Response, error) {
	id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	if id.IsValid() && id.CanSet() && id.String() == "" {
		id.SetString(genId())
	}
	return c.client.Upsert(coll+"docs/", doc, &doc, opts...)
}

// Delete database
func (c *CosmosDB) DeleteDatabase(link string) (*Response, error) {
	return c.client.Delete(link)
}

// Delete collection
func (c *CosmosDB) DeleteCollection(link string) (*Response, error) {
	return c.client.Delete(link)
}

// Delete collection
func (c *CosmosDB) DeleteDocument(link string) (*Response, error) {
	return c.client.Delete(link)
}

// Delete stored procedure
func (c *CosmosDB) DeleteStoredProcedure(link string) (*Response, error) {
	return c.client.Delete(link)
}

// Delete user defined function
func (c *CosmosDB) DeleteUserDefinedFunction(link string) (*Response, error) {
	return c.client.Delete(link)
}

// Replace database
func (c *CosmosDB) ReplaceDatabase(link string, body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Replace(link, body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Replace document
func (c *CosmosDB) ReplaceDocument(link string, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.Replace(link, doc, &doc, opts...)
}

// Replace Async document
func (c *CosmosDB) ReplaceDocumentAsync(link string, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.ReplaceAsync(link, doc, &doc, opts...)
}

// Replace stored procedure
func (c *CosmosDB) ReplaceStoredProcedure(link string, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Replace(link, body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Replace stored procedure
func (c *CosmosDB) ReplaceUserDefinedFunction(link string, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Replace(link, body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Execute stored procedure
func (c *CosmosDB) ExecuteStoredProcedure(link string, params, body interface{}, opts ...CallOption) (err error) {
	_, err = c.client.Execute(link, params, &body, opts...)
	return
}

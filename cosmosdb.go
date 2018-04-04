//
// This library started as a fork of `github.com/nerdylikeme/go-documentdb`
//

package gocosmosdb

import "reflect"

type Config struct {
	MasterKey string
}

type CosmosDB struct {
	client Clienter
}

// Create CosmosDBClient
func New(url string, config Config) *CosmosDB {
	client := &Client{}
	client.Url = url
	client.Config = config
	return &CosmosDB{client}
}

// TODO: Add `requestOptions` arguments
// Read database by self link
func (c *CosmosDB) ReadDatabase(link string) (db *Database, err error) {
	err = c.client.Read(link, &db)
	if err != nil {
		return nil, err
	}
	return
}

// Read collection by self link
func (c *CosmosDB) ReadCollection(link string) (coll *Collection, err error) {
	err = c.client.Read(link, &coll)
	if err != nil {
		return nil, err
	}
	return
}

// Read document by self link
func (c *CosmosDB) ReadDocument(link string, doc interface{}) (err error) {
	err = c.client.Read(link, &doc)
	return
}

// Read sporc by self link
func (c *CosmosDB) ReadStoredProcedure(link string) (sproc *Sproc, err error) {
	err = c.client.Read(link, &sproc)
	if err != nil {
		return nil, err
	}
	return
}

// Read udf by self link
func (c *CosmosDB) ReadUserDefinedFunction(link string) (udf *UDF, err error) {
	err = c.client.Read(link, &udf)
	if err != nil {
		return nil, err
	}
	return
}

// Read all databases
func (c *CosmosDB) ReadDatabases() (dbs []Database, err error) {
	return c.QueryDatabases("")
}

// Read all collections by db selflink
func (c *CosmosDB) ReadCollections(db string) (colls []Collection, err error) {
	return c.QueryCollections(db, "")
}

// Read all sprocs by collection self link
func (c *CosmosDB) ReadStoredProcedures(coll string) (sprocs []Sproc, err error) {
	return c.QueryStoredProcedures(coll, "")
}

// Read all udfs by collection self link
func (c *CosmosDB) ReadUserDefinedFunctions(coll string) (udfs []UDF, err error) {
	return c.QueryUserDefinedFunctions(coll, "")
}

// Read all collection documents by self link
// TODO: use iterator for heavy transactions
func (c *CosmosDB) ReadDocuments(coll string, docs interface{}) (err error) {
	return c.QueryDocuments(coll, "", docs)
}

// Read all databases that satisfy a query
func (c *CosmosDB) QueryDatabases(query string) (dbs []Database, err error) {
	data := struct {
		Databases []Database `json:"Databases,omitempty"`
		Count     int        `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		err = c.client.Query("dbs", query, &data)
	} else {
		err = c.client.Read("dbs", &data)
	}
	if dbs = data.Databases; err != nil {
		dbs = nil
	}
	return
}

// Read all db-collection that satisfy a query
func (c *CosmosDB) QueryCollections(db, query string) (colls []Collection, err error) {
	data := struct {
		Collections []Collection `json:"DocumentCollections,omitempty"`
		Count       int          `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		err = c.client.Query(db+"colls/", query, &data)
	} else {
		err = c.client.Read(db+"colls/", &data)
	}
	if colls = data.Collections; err != nil {
		colls = nil
	}
	return
}

// Read all collection `sprocs` that satisfy a query
func (c *CosmosDB) QueryStoredProcedures(coll, query string) (sprocs []Sproc, err error) {
	data := struct {
		Sprocs []Sproc `json:"StoredProcedures,omitempty"`
		Count  int     `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		err = c.client.Query(coll+"sprocs/", query, &data)
	} else {
		err = c.client.Read(coll+"sprocs/", &data)
	}
	if sprocs = data.Sprocs; err != nil {
		sprocs = nil
	}
	return
}

// Read all collection `udfs` that satisfy a query
func (c *CosmosDB) QueryUserDefinedFunctions(coll, query string) (udfs []UDF, err error) {
	data := struct {
		Udfs  []UDF `json:"UserDefinedFunctions,omitempty"`
		Count int   `json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		err = c.client.Query(coll+"udfs/", query, &data)
	} else {
		err = c.client.Read(coll+"udfs/", &data)
	}
	if udfs = data.Udfs; err != nil {
		udfs = nil
	}
	return
}

// Read all documents in a collection that satisfy a query
func (c *CosmosDB) QueryDocuments(coll, query string, docs interface{}) (err error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if len(query) > 0 {
		err = c.client.Query(coll+"docs/", query, &data)
	} else {
		err = c.client.Read(coll+"docs/", &data)
	}
	return
}

// Create database
func (c *CosmosDB) CreateDatabase(body interface{}) (db *Database, err error) {
	err = c.client.Create("dbs", body, &db)
	if err != nil {
		return nil, err
	}
	return
}

// Create collection
func (c *CosmosDB) CreateCollection(db string, body interface{}) (coll *Collection, err error) {
	err = c.client.Create(db+"colls/", body, &coll)
	if err != nil {
		return nil, err
	}
	return
}

// Create stored procedure
func (c *CosmosDB) CreateStoredProcedure(coll string, body interface{}) (sproc *Sproc, err error) {
	err = c.client.Create(coll+"sprocs/", body, &sproc)
	if err != nil {
		return nil, err
	}
	return
}

// Create user defined function
func (c *CosmosDB) CreateUserDefinedFunction(coll string, body interface{}) (udf *UDF, err error) {
	err = c.client.Create(coll+"udfs/", body, &udf)
	if err != nil {
		return nil, err
	}
	return
}

// Create document
func (c *CosmosDB) CreateDocument(coll string, doc interface{}) error {
	id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	if id.IsValid() && id.String() == "" {
		id.SetString(uuid())
	}
	return c.client.Create(coll+"docs/", doc, &doc)
}

// TODO: DRY, but the sdk want that[mm.. maybe just client.Delete(self_link)]
// Delete database
func (c *CosmosDB) DeleteDatabase(link string) error {
	return c.client.Delete(link)
}

// Delete collection
func (c *CosmosDB) DeleteCollection(link string) error {
	return c.client.Delete(link)
}

// Delete collection
func (c *CosmosDB) DeleteDocument(link string) error {
	return c.client.Delete(link)
}

// Delete stored procedure
func (c *CosmosDB) DeleteStoredProcedure(link string) error {
	return c.client.Delete(link)
}

// Delete user defined function
func (c *CosmosDB) DeleteUserDefinedFunction(link string) error {
	return c.client.Delete(link)
}

// Replace database
func (c *CosmosDB) ReplaceDatabase(link string, body interface{}) (db *Database, err error) {
	err = c.client.Replace(link, body, &db)
	if err != nil {
		return nil, err
	}
	return
}

// Replace document
func (c *CosmosDB) ReplaceDocument(link string, doc interface{}) error {
	return c.client.Replace(link, doc, &doc)
}

// Replace stored procedure
func (c *CosmosDB) ReplaceStoredProcedure(link string, body interface{}) (sproc *Sproc, err error) {
	err = c.client.Replace(link, body, &sproc)
	if err != nil {
		return nil, err
	}
	return
}

// Replace stored procedure
func (c *CosmosDB) ReplaceUserDefinedFunction(link string, body interface{}) (udf *UDF, err error) {
	err = c.client.Replace(link, body, &udf)
	if err != nil {
		return nil, err
	}
	return
}

// Execute stored procedure
func (c *CosmosDB) ExecuteStoredProcedure(link string, params, body interface{}) (err error) {
	err = c.client.Execute(link, params, &body)
	return
}
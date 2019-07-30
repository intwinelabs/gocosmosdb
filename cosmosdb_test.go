package gocosmosdb

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// for advanced debugging
// client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", Debug: true, Verbose: true}, log)
func TestNew(t *testing.T) {
	assert := assert.New(t)
	client := New("url", Config{MasterKey: "config"}, log)
	assert.IsType(client, &CosmosDB{}, "Should return CosmosDB object")
	conf := client.GetConfig()
	assert.Equal(Config{MasterKey: "config"}, conf)
	uri := client.GetURI()
	assert.Equal("url", uri)
	client.EnableDebug()
	conf = client.GetConfig()
	assert.Equal(true, conf.Debug)
	client.DisableDebug()
	conf = client.GetConfig()
	assert.Equal(false, conf.Debug)
}

func TestNewRetryable(t *testing.T) {
	assert := assert.New(t)
	c := Config{
		MasterKey:    "config",
		RetryWaitMin: 100 * time.Millisecond,
		RetryWaitMax: 100 * time.Millisecond,
		RetryMax:     2,
		Pooled:       true,
	}
	client := New("url", c, log)
	assert.IsType(client, &CosmosDB{}, "Should return CosmosDB object")
	conf := client.GetConfig()
	assert.Equal(conf.RetryWaitMin, 100*time.Millisecond)
	assert.Equal(client.client.httpClient.RetryWaitMin, 100*time.Millisecond)
	assert.Equal(conf.RetryWaitMax, 100*time.Millisecond)
	assert.Equal(client.client.httpClient.RetryWaitMax, 100*time.Millisecond)
	assert.Equal(conf.RetryMax, 2)
	assert.Equal(client.client.httpClient.RetryMax, 2)
	assert.Equal(conf.Pooled, true)
	assert.NotNil(client.client.httpClient.HTTPClient.Transport)
}

func TestReadDatabase(t *testing.T) {
	assert := assert.New(t)
	resp := `{
		"id": "iot2",
		"_rid": "qicAAA==",
		"_ts": 1446192371,
		"_self": "dbs\/qicAAA==\/",
		"_etag": "\"00001800-0000-0000-0000-563324f30000\"",
		"_colls": "colls\/",
		"_users": "users\/"
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	db, err := client.ReadDatabase("dbs/qicAAA==")
	assert.Nil(err)
	assert.Equal("iot2", db.Id)
}

func TestReadDatabaseWithDebugging(t *testing.T) {
	assert := assert.New(t)
	resp := `{
		"id": "iot2",
		"_rid": "qicAAA==",
		"_ts": 1446192371,
		"_self": "dbs\/qicAAA==\/",
		"_etag": "\"00001800-0000-0000-0000-563324f30000\"",
		"_colls": "colls\/",
		"_users": "users\/"
	}`
	s := ServerFactory(resp)
	s.SetStatus(222)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", Debug: true, Verbose: true}, log)
	_, err := client.ReadDatabase("dbs/qicAAA==")
	assert.NotNil(err)
}

func TestReadCollection(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SampleCollection",  
		"indexingPolicy": {  
		  "indexingMode": "consistent",  
		  "automatic": true,  
		  "includedPaths": [  
			{  
			  "path": "/*",  
			  "indexes": [  
				{  
				  "kind": "Range",  
				  "dataType": "Number",  
				  "precision": -1  
				},  
				{  
				  "kind": "Hash",  
				  "dataType": "String",  
				  "precision": 3  
				}  
			  ]  
			}  
		  ],  
		  "excludedPaths": []  
		},  
		"_rid": "PaYSAPH7qAo=",  
		"_ts": 1459194239,  
		"_self": "dbs/PaYSAA==/colls/PaYSAPH7qAo=/",  
		"_etag": "\"00001300-0000-0000-0000-56f9897f0000\"",  
		"_docs": "docs/",  
		"_sprocs": "sprocs/",  
		"_triggers": "triggers/",  
		"_udfs": "udfs/",  
		"_conflicts": "conflicts/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	coll, err := client.ReadCollection("dbs/PaYSAA==/colls/PaYSAPH7qAo=")
	assert.Nil(err)
	assert.Equal("SampleCollection", coll.Id)
}

type testDoc struct {
	Document
	PONumber string `json:"ponumber"`
}

func TestReadDocument(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	doc := testDoc{}
	err := client.ReadDocument("dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==", &doc)
	assert.Nil(err)
	assert.Equal("SalesOrder1", doc.Id)
}

func TestReadStoredProcedure(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World!\");\r\n}",  
		"id": "sproc_hello_world",  
		"_rid": "Sl8fALN4sw4CAAAAAAAAgA==",  
		"_ts": 1449681197,  
		"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==/",  
		"_etag": "\"06003ce1-0000-0000-0000-5668612d0000\""  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	sproc, err := client.ReadStoredProcedure("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==")
	assert.Nil(err)
	assert.Equal("sproc_hello_world", sproc.Id)
}

func TestReadUserDefinedFunction(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
        "body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
        "id": "simpleTaxUDF",  
        "_rid": "Sl8fALN4sw4BAAAAAAAAYA==",  
        "_ts": 1449688293,  
        "_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs/Sl8fALN4sw4BAAAAAAAAYA==/",  
        "_etag": "\"060072e4-0000-0000-0000-56687ce50000\""  
    }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	udf, err := client.ReadUserDefinedFunction("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs/Sl8fALN4sw4BAAAAAAAAYA==")
	assert.Nil(err)
	assert.Equal("simpleTaxUDF", udf.Id)
}

func TestReadDatabases(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "",  
		"Databases": [{  
			"id": "iot2",  
			"_rid": "qicAAA==",  
			"_ts": 1446192371,  
			"_self": "dbs\/qicAAA==\/",  
			"_etag": "\"00001800-0000-0000-0000-563324f30000\"",  
			"_colls": "colls\/",  
			"_users": "users\/"  
		},  
		{  
			"id": "TestDB2",  
			"_rid": "KI0YAA==",  
			"_ts": 1446243863,  
			"_self": "dbs\/KI0YAA==\/",  
			"_etag": "\"00001f00-0000-0000-0000-5633ee170000\"",  
			"_colls": "colls\/",  
			"_users": "users\/"  
		},  
		{  
			"id": "FoodDB",  
			"_rid": "vdoeAA==",  
			"_ts": 1442511602,  
			"_self": "dbs\/vdoeAA==\/",  
			"_etag": "\"00000100-0000-0000-0000-55fafaf20000\"",  
			"_colls": "colls\/",  
			"_users": "users\/"  
		}],  
		"_count": 3  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	dbs, err := client.ReadDatabases()
	assert.Nil(err)
	assert.Equal("iot2", dbs[0].Id)
	assert.Equal("TestDB2", dbs[1].Id)
	assert.Equal("FoodDB", dbs[2].Id)
}

func TestReadCollections(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "PaYSAA==",  
		"DocumentCollections": [  
		  {  
			"id": "SampleCollection",  
			"indexingPolicy": {  
			  "indexingMode": "consistent",  
			  "automatic": true,  
			  "includedPaths": [  
				{  
				  "path": "/*",  
				  "indexes": [  
					{  
					  "kind": "Range",  
					  "dataType": "Number",  
					  "precision": -1  
					},  
					{  
					  "kind": "Hash",  
					  "dataType": "String",  
					  "precision": 3  
					}  
				  ]  
				}  
			  ],  
			  "excludedPaths": []  
			},  
			"_rid": "PaYSAPH7qAo=",  
			"_ts": 1459194239,  
			"_self": "dbs/PaYSAA==/colls/PaYSAPH7qAo=/",  
			"_etag": "\"00001300-0000-0000-0000-56f9897f0000\"",  
			"_docs": "docs/",  
			"_sprocs": "sprocs/",  
			"_triggers": "triggers/",  
			"_udfs": "udfs/",  
			"_conflicts": "conflicts/"  
		  },  
		  {  
			"id": "SampleCollectionWithCustomIndexPolicy",  
			"indexingPolicy": {  
			  "indexingMode": "lazy",  
			  "automatic": true,  
			  "includedPaths": [  
				{  
				  "path": "/*",  
				  "indexes": [  
					{  
					  "kind": "Range",  
					  "dataType": "Number",  
					  "precision": -1  
					},  
					{  
					  "kind": "Hash",  
					  "dataType": "String",  
					  "precision": 3  
					}  
				  ]  
				}  
			  ],  
			  "excludedPaths": []  
			},  
			"_rid": "PaYSAIxUPws=",  
			"_ts": 1459194241,  
			"_self": "dbs/PaYSAA==/colls/PaYSAIxUPws=/",  
			"_etag": "\"00001500-0000-0000-0000-56f989810000\"",  
			"_docs": "docs/",  
			"_sprocs": "sprocs/",  
			"_triggers": "triggers/",  
			"_udfs": "udfs/",  
			"_conflicts": "conflicts/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	colls, err := client.ReadCollections("dbs/PaYSAA==")
	assert.Nil(err)
	assert.Equal("SampleCollection", colls[0].Id)
	assert.Equal("SampleCollectionWithCustomIndexPolicy", colls[1].Id)
}

func TestReadDocuments(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "d9RzAJRFKgw=",  
		"Documents": [  
		  {  
			"id": "SalesOrder1",  
			"ponumber": "PO18009186470",  
			"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
			"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  },  
		  {  
			"id": "SalesOrder2",  
			"ponumber": "PO15428132599",  
			"_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
			"_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	docs := []testDoc{}
	err := client.ReadDocuments("dbs/d9RzAA==/colls/d9RzAJRFKgw=", &docs)
	assert.Nil(err)
	assert.Equal("SalesOrder1", docs[0].Id)
	assert.Equal("SalesOrder2", docs[1].Id)
}

func TestReadStoredProcedures(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "Sl8fALN4sw4=",  
		"StoredProcedures": [{  
			"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World!\");\r\n}",  
			"id": "sproc_hello_world",  
			"_rid": "Sl8fALN4sw4CAAAAAAAAgA==",  
			"_ts": 1449681197,  
			"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==/",  
			"_etag": "\"06003ce1-0000-0000-0000-5668612d0000\""  
		}],  
		"_count": 1  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	sprocs, err := client.ReadStoredProcedures("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs")
	assert.Nil(err)
	assert.Equal("sproc_hello_world", sprocs[0].Id)
}

func TestReadUserDefinedFunctions(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "Sl8fALN4sw4=",  
		"UserDefinedFunctions": [{  
			"body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
			"id": "simpleTaxUDF",  
			"_rid": "Sl8fALN4sw4BAAAAAAAAYA==",  
			"_ts": 1449688293,  
			"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs/Sl8fALN4sw4BAAAAAAAAYA==/",  
			"_etag": "\"060072e4-0000-0000-0000-56687ce50000\""  
		}],  
		"_count": 1  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	udfs, err := client.ReadUserDefinedFunctions("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs")
	assert.Nil(err)
	assert.Equal("simpleTaxUDF", udfs[0].Id)
}

func TestQueryDatabases(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "",  
		"Databases": [{  
			"id": "iot2",  
			"_rid": "qicAAA==",  
			"_ts": 1446192371,  
			"_self": "dbs\/qicAAA==\/",  
			"_etag": "\"00001800-0000-0000-0000-563324f30000\"",  
			"_colls": "colls\/",  
			"_users": "users\/"  
		},  
		{  
			"id": "TestDB2",  
			"_rid": "KI0YAA==",  
			"_ts": 1446243863,  
			"_self": "dbs\/KI0YAA==\/",  
			"_etag": "\"00001f00-0000-0000-0000-5633ee170000\"",  
			"_colls": "colls\/",  
			"_users": "users\/"  
		},  
		{  
			"id": "FoodDB",  
			"_rid": "vdoeAA==",  
			"_ts": 1442511602,  
			"_self": "dbs\/vdoeAA==\/",  
			"_etag": "\"00000100-0000-0000-0000-55fafaf20000\"",  
			"_colls": "colls\/",  
			"_users": "users\/"  
		}],  
		"_count": 3  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	dbs, err := client.QueryDatabases("SELECT * FROM root")
	assert.Nil(err)
	assert.Equal("iot2", dbs[0].Id)
	assert.Equal("TestDB2", dbs[1].Id)
	assert.Equal("FoodDB", dbs[2].Id)
}

func TestQueryCollections(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "PaYSAA==",  
		"DocumentCollections": [  
		  {  
			"id": "SampleCollection",  
			"indexingPolicy": {  
			  "indexingMode": "consistent",  
			  "automatic": true,  
			  "includedPaths": [  
				{  
				  "path": "/*",  
				  "indexes": [  
					{  
					  "kind": "Range",  
					  "dataType": "Number",  
					  "precision": -1  
					},  
					{  
					  "kind": "Hash",  
					  "dataType": "String",  
					  "precision": 3  
					}  
				  ]  
				}  
			  ],  
			  "excludedPaths": []  
			},  
			"_rid": "PaYSAPH7qAo=",  
			"_ts": 1459194239,  
			"_self": "dbs/PaYSAA==/colls/PaYSAPH7qAo=/",  
			"_etag": "\"00001300-0000-0000-0000-56f9897f0000\"",  
			"_docs": "docs/",  
			"_sprocs": "sprocs/",  
			"_triggers": "triggers/",  
			"_udfs": "udfs/",  
			"_conflicts": "conflicts/"  
		  },  
		  {  
			"id": "SampleCollectionWithCustomIndexPolicy",  
			"indexingPolicy": {  
			  "indexingMode": "lazy",  
			  "automatic": true,  
			  "includedPaths": [  
				{  
				  "path": "/*",  
				  "indexes": [  
					{  
					  "kind": "Range",  
					  "dataType": "Number",  
					  "precision": -1  
					},  
					{  
					  "kind": "Hash",  
					  "dataType": "String",  
					  "precision": 3  
					}  
				  ]  
				}  
			  ],  
			  "excludedPaths": []  
			},  
			"_rid": "PaYSAIxUPws=",  
			"_ts": 1459194241,  
			"_self": "dbs/PaYSAA==/colls/PaYSAIxUPws=/",  
			"_etag": "\"00001500-0000-0000-0000-56f989810000\"",  
			"_docs": "docs/",  
			"_sprocs": "sprocs/",  
			"_triggers": "triggers/",  
			"_udfs": "udfs/",  
			"_conflicts": "conflicts/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	colls, err := client.QueryCollections("dbs/PaYSAA==", "SELECT * FROM root")
	assert.Nil(err)
	assert.Equal("SampleCollection", colls[0].Id)
	assert.Equal("SampleCollectionWithCustomIndexPolicy", colls[1].Id)
}

func TestQueryStoredProcedures(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "Sl8fALN4sw4=",  
		"StoredProcedures": [{  
			"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World!\");\r\n}",  
			"id": "sproc_hello_world",  
			"_rid": "Sl8fALN4sw4CAAAAAAAAgA==",  
			"_ts": 1449681197,  
			"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==/",  
			"_etag": "\"06003ce1-0000-0000-0000-5668612d0000\""  
		}],  
		"_count": 1  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	sprocs, err := client.QueryStoredProcedures("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs", "SELECT * FROM root")
	assert.Nil(err)
	assert.Equal("sproc_hello_world", sprocs[0].Id)
}

func TestQueryUserDefinedFunctions(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "Sl8fALN4sw4=",  
		"UserDefinedFunctions": [{  
			"body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
			"id": "simpleTaxUDF",  
			"_rid": "Sl8fALN4sw4BAAAAAAAAYA==",  
			"_ts": 1449688293,  
			"_self": "dbs\/Sl8fAA==\/colls\/Sl8fALN4sw4=\/udfs\/Sl8fALN4sw4BAAAAAAAAYA==\/",  
			"_etag": "\"060072e4-0000-0000-0000-56687ce50000\""  
		}],  
		"_count": 1  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	udfs, err := client.QueryUserDefinedFunctions("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs", "SELECT * FROM root")
	assert.Nil(err)
	assert.Equal("simpleTaxUDF", udfs[0].Id)
}

func TestQueryDocuments(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "d9RzAJRFKgw=",  
		"Documents": [  
		  {  
			"id": "SalesOrder1",  
			"ponumber": "PO18009186470",  
			"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
			"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  },  
		  {  
			"id": "SalesOrder2",  
			"ponumber": "PO15428132599",  
			"_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
			"_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	docs := []testDoc{}
	err := client.QueryDocuments("dbs/d9RzAA==/colls/d9RzAJRFKgw=", "SELECT * FROM root", &docs)
	assert.Nil(err)
	assert.Equal("SalesOrder1", docs[0].Id)
	assert.Equal("SalesOrder2", docs[1].Id)
}

func TestQueryDocumentsBadQuery(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(400)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", Verbose: true, Debug: true}, log)
	docs := []testDoc{}
	err := client.QueryDocuments("dbs/d9RzAA==/colls/d9RzAJRFKgw=", "SELECT * root", &docs)
	assert.NotNil(err)
}

func TestQueryDocumentsWithPartitionKey(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "d9RzAJRFKgw=",  
		"Documents": [  
		  {  
			"id": "SalesOrder1",  
			"ponumber": "PO18009186470",  
			"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
			"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  },  
		  {  
			"id": "SalesOrder2",  
			"ponumber": "PO15428132599",  
			"_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
			"_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", PartitionKeyStructField: "PONumber", PartitionKeyPath: "/ponumber"}, log)
	docs := []testDoc{}
	err := client.QueryDocuments("dbs/d9RzAA==/colls/d9RzAJRFKgw=", "SELECT * FROM root", &docs)
	assert.Nil(err)
	assert.Equal("SalesOrder1", docs[0].Id)
	assert.Equal("SalesOrder2", docs[1].Id)
}

func TestQueryDocumentsWithParameters(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "d9RzAJRFKgw=",  
		"Documents": [  
		  {  
			"id": "SalesOrder1",  
			"ponumber": "PO18009186470",  
			"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
			"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  },  
		  {  
			"id": "SalesOrder2",  
			"ponumber": "PO15428132599",  
			"_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
			"_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	docs := []testDoc{}
	query := &QueryWithParameters{
		Query: "SELECT * FROM root r WHERE r._ts > @_ts",
		Parameters: []QueryParameter{
			QueryParameter{
				Name:  "@_ts",
				Value: 1459216957,
			},
		},
	}
	err := client.QueryDocumentsWithParameters("dbs/d9RzAA==/colls/d9RzAJRFKgw=", query, &docs)
	assert.Nil(err)
	assert.Equal("SalesOrder1", docs[0].Id)
	assert.Equal("SalesOrder2", docs[1].Id)
}

func TestQueryDocumentsWithParametersWithPartitionKey(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"_rid": "d9RzAJRFKgw=",  
		"Documents": [  
		  {  
			"id": "SalesOrder1",  
			"ponumber": "PO18009186470",  
			"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
			"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  },  
		  {  
			"id": "SalesOrder2",  
			"ponumber": "PO15428132599",  
			"_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
			"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
			"_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
			"_ts": 1459216987,  
			"_attachments": "attachments/"  
		  }  
		],  
		"_count": 2  
	  }`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", PartitionKeyStructField: "PONumber", PartitionKeyPath: "/ponumber"}, log)
	docs := []testDoc{}
	query := &QueryWithParameters{
		Query: "SELECT * FROM root r WHERE r._ts > @_ts",
		Parameters: []QueryParameter{
			QueryParameter{
				Name:  "@_ts",
				Value: 1459216957,
			},
		},
	}
	err := client.QueryDocumentsWithParameters("dbs/d9RzAA==/colls/d9RzAJRFKgw=", query, &docs)
	assert.Nil(err)
	assert.Equal("SalesOrder1", docs[0].Id)
	assert.Equal("SalesOrder2", docs[1].Id)
}

func TestQueryPartitionKeyRanges(t *testing.T) {
	assert := assert.New(t)
	resp := `{
		"_rid":"qYcAAPEvJBQ=",
		"PartitionKeyRanges":[
		   {
			  "_rid":"qYcAAPEvJBQCAAAAAAAAUA==",
			  "id":"0",
			  "_etag":"\"00002800-0000-0000-0000-580ac4ea0000\"",
			  "minInclusive":"",
			  "maxExclusive":"05C1CFFFFFFFF8",
			  "_self":"dbs/qYcAAA==/colls/qYcAAPEvJBQ=/pkranges/qYcAAPEvJBQCAAAAAAAAUA==/",
			  "_ts":1477100776
		   }
		],
		"_count": 1
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	ranges, err := client.QueryPartitionKeyRanges("dbs/qYcAAA==/colls/qYcAAPEvJBQ=/pkranges", "SELECT * FROM root")
	assert.Nil(err)
	assert.Equal("0", ranges[0].Id)
}

func TestCreateDatabase(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "volcanodb2",  
		"_rid": "CqNBAA==",  
		"_ts": 1449602962,  
		"_self": "dbs\/CqNBAA==\/",  
		"_etag": "\"00000a00-0000-0000-0000-56672f920000\"",  
		"_colls": "colls/",  
		"_users": "users/"  
	}`
	s := ServerFactory(resp)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	db, err := client.CreateDatabase(`{"id": "volcanodb2"}`)
	t.Logf("%+v", db)
	assert.Nil(err)
	//assert.Equal("volcanodb2", db.Id)
}

func TestCreateCollection(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "testcoll",  
		"indexingPolicy": {  
		  "automatic": true,  
		  "indexingMode": "Consistent",  
		  "includedPaths": [  
			{  
			  "path": "/*",  
			  "indexes": [  
				{  
				  "dataType": "String",  
				  "precision": -1,  
				  "kind": "Range"  
				}  
			  ]  
			}  
		  ]  
		},  
		"partitionKey": {  
		  "paths": [  
			"/AccountNumber"  
		  ],  
		  "kind": "Hash",
		   "Version": 2
	  
		}  
	  }`
	s := ServerFactory(resp)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	coll, err := client.CreateCollection("dbs/qYcAAA==", `{"id": "testcoll"}`)
	assert.Nil(err)
	assert.Equal("testcoll", coll.Id)
}

func TestCreateStoredProcedure(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World\");\r\n}",  
		"id": "sproc_1",  
		"_rid": "Sl8fALN4sw4CAAAAAAAAgA==",  
		"_ts": 1449680569,  
		"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==/",  
		"_etag": "\"0600ffe0-0000-0000-0000-56685eb90000\""  
	}`
	s := ServerFactory(resp)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	sprocBody := `{  
		"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World\");\r\n}",  
		"id": "sproc_1"  
	}`
	sproc, err := client.CreateStoredProcedure("dbs/qYcAAA==/colls/qYcAAPEvJBQ=/sprocs", sprocBody)
	assert.Nil(err)
	assert.Equal("sproc_1", sproc.Id)
}

func TestCreateUserDefinedFunction(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
		"id": "simpleTaxUDF",  
		"_rid": "Sl8fALN4sw4BAAAAAAAAYA==",  
		"_ts": 1449687949,  
		"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs/Sl8fALN4sw4BAAAAAAAAYA==/",  
		"_etag": "\"06003ee4-0000-0000-0000-56687b8d0000\""  
	}`
	s := ServerFactory(resp)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	udfBody := `{  
		"body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
		"id": "simpleTaxUDF"  
	}`
	udf, err := client.CreateUserDefinedFunction("dbs/qYcAAA==/colls/qYcAAPEvJBQ=/udfs", udfBody)
	assert.Nil(err)
	assert.Equal("simpleTaxUDF", udf.Id)
}

func TestCreateDocument(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	r, err := client.CreateDocument("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestCreateDocumentWithPartitionKey(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", PartitionKeyStructField: "PONumber", PartitionKeyPath: "/ponumber"}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	r, err := client.CreateDocument("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestUpsertDocument(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	r, err := client.UpsertDocument("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestUpsertDocumentWithPartitionKey(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", PartitionKeyStructField: "PONumber", PartitionKeyPath: "/ponumber"}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	r, err := client.UpsertDocument("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestDeleteDatabase(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(204)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	_, err := client.DeleteDatabase("dbs/qYcAAA==")
	assert.Nil(err)
}

func TestDeleteCollection(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(204)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	_, err := client.DeleteCollection("dbs/qYcAAA==/colls/qYcAAPEvJBQ=")
	assert.Nil(err)
}

func TestDeleteDocument(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(204)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	_, err := client.DeleteDocument("dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==")
	assert.Nil(err)
}

func TestDeleteStoredProcedure(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(204)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	_, err := client.DeleteStoredProcedure("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==")
	assert.Nil(err)
}

func TestDeleteUserDefinedFunctions(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(204)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	_, err := client.DeleteUserDefinedFunction("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs/Sl8fALN4sw4BAAAAAAAAYA==")
	assert.Nil(err)
}

func TestReplaceDatabase(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "newid",  
		"_rid": "CqNBAA==",  
		"_ts": 1449602962,  
		"_self": "dbs\/CqNBAA==\/",  
		"_etag": "\"00000a00-0000-0000-0000-56672f920000\"",  
		"_colls": "colls\/",  
		"_users": "users\/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	db, err := client.ReplaceDatabase("dbs/qYcAAA==", `{"id": "newid"}`)
	assert.Nil(err)
	assert.Equal("newid", db.Id)
}

func TestReplaceDocument(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	r, err := client.ReplaceDocument("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestReplaceDocumentWithPartitionKey(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", PartitionKeyStructField: "PONumber", PartitionKeyPath: "/ponumber"}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	r, err := client.ReplaceDocument("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestReplaceDocumentAsync(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	doc.Etag = "\"0000d986-0000-0000-0000-56f9e25b0000\""
	r, err := client.ReplaceDocumentAsync("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestReplaceDocumentAsyncWithPartitionKey(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"id": "SalesOrder1",  
		"ponumber": "PO18009186470",  
		"_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		"_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		"_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		"_ts": 1459216987,  
		"_attachments": "attachments/"  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", PartitionKeyStructField: "PONumber", PartitionKeyPath: "/ponumber"}, log)
	doc := testDoc{}
	doc.Id = "SalesOrder1"
	doc.PONumber = "PO18009186470"
	doc.Etag = "\"0000d986-0000-0000-0000-56f9e25b0000\""
	r, err := client.ReplaceDocumentAsync("dbs/qYcAAA==/colls/qYcAAPEvJBQ=", &doc)
	assert.Nil(err)
	assert.NotNil(r)
}

func TestRepalceStoredProcedure(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World\");\r\n}",  
		"id": "new_sproc_1",  
		"_rid": "Sl8fALN4sw4CAAAAAAAAgA==",  
		"_ts": 1449680569,  
		"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==/",  
		"_etag": "\"0600ffe0-0000-0000-0000-56685eb90000\""  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	sprocBody := `{  
		"body": "function () {\r\n    var context = getContext();\r\n    var response = context.getResponse();\r\n\r\n    response.setBody(\"Hello, World\");\r\n}",  
		"id": "new_sproc_1"  
	}`
	sproc, err := client.ReplaceStoredProcedure("dbs/qYcAAA==/colls/qYcAAPEvJBQ=/sprocs", sprocBody)
	assert.Nil(err)
	assert.Equal("new_sproc_1", sproc.Id)
}
func TestReplaceUserDefinedFunction(t *testing.T) {
	assert := assert.New(t)
	resp := `{  
		"body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
		"id": "newSimpleTaxUDF",  
		"_rid": "Sl8fALN4sw4BAAAAAAAAYA==",  
		"_ts": 1449687949,  
		"_self": "dbs/Sl8fAA==/colls/Sl8fALN4sw4=/udfs/Sl8fALN4sw4BAAAAAAAAYA==/",  
		"_etag": "\"06003ee4-0000-0000-0000-56687b8d0000\""  
	}`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	udfBody := `{  
		"body": "function tax(income) {\r\n    if(income == undefined) \r\n        throw 'no input';\r\n    if (income < 1000) \r\n        return income * 0.1;\r\n    else if (income < 10000) \r\n        return income * 0.2;\r\n    else\r\n        return income * 0.4;\r\n}",  
		"id": "newSimpleTaxUDF"  
	}`
	udf, err := client.ReplaceUserDefinedFunction("dbs/qYcAAA==/colls/qYcAAPEvJBQ=/udfs", udfBody)
	assert.Nil(err)
	assert.Equal("newSimpleTaxUDF", udf.Id)
}

func TestExecuteStoredProcedure(t *testing.T) {
	assert := assert.New(t)
	resp := `[  
		{  
		  "id": "SalesOrder1",  
		  "ponumber": "PO18009186470",  
		  "_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
		  "_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
		  "_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
		  "_ts": 1459216987,  
		  "_attachments": "attachments/"  
		},  
		{  
		  "id": "SalesOrder2",  
		  "ponumber": "PO15428132599",  
		  "_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
		  "_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
		  "_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
		  "_ts": 1459216987,  
		  "_attachments": "attachments/"  
		}  
	  ]`
	s := ServerFactory(resp)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	docs := []testDoc{}
	err := client.ExecuteStoredProcedure("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==", []string{"param1"}, &docs)
	assert.Nil(err)
	assert.Equal("SalesOrder1", docs[0].Id)
	assert.Equal("SalesOrder2", docs[1].Id)
}

func TestExecuteStoredProcedureWithContextCancel(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(500, 500, 500, 500)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", RetryWaitMin: 100 * time.Millisecond, RetryWaitMax: 1 * time.Millisecond, RetryMax: 3}, log)
	ctx, cancel := context.WithCancel(context.Background())
	docs := []testDoc{}
	go func() {
		time.Sleep(250 * time.Microsecond)
		cancel()
	}()
	err := client.ExecuteStoredProcedure("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==", []string{"param1"}, &docs, WithContext(ctx))
	assert.NotNil(err)
	assert.Contains(err.Error(), "context canceled")

}

func TestExecuteStoredProcedureWithContextTimeout(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(500, 500, 500, 500)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", RetryWaitMin: 100 * time.Millisecond, RetryWaitMax: 100 * time.Millisecond, RetryMax: 3}, log)
	ctx, _ := context.WithTimeout(context.Background(), 250*time.Millisecond)
	docs := []testDoc{}
	err := client.ExecuteStoredProcedure("dbs/Sl8fAA==/colls/Sl8fALN4sw4=/sprocs/Sl8fALN4sw4CAAAAAAAAgA==", []string{"param1"}, &docs, WithContext(ctx))
	assert.NotNil(err)
	assert.Contains(err.Error(), "context deadline exceeded")

}

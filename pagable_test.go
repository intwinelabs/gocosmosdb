package gocosmosdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagable(t *testing.T) {
	assert := assert.New(t)
	resp1 := `{  
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
		  }  
		],  
		"_count": 2  
	  }`
	resp2 := `{  
		"_rid": "d9RzAJRFKgw=",  
		"Documents": [  
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
	s := ServerFactory(resp1, resp2)
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
	pg := client.NewPagableQuery("dbs/d9RzAA==/colls/d9RzAJRFKgw=", query, 1, &docs)
	assert.NotNil(pg)
	assert.IsType(&PagableQuery{}, pg)
	pg.Next()
	assert.Equal("SalesOrder1", docs[0].Id)
	pg.Next()
	assert.Equal("SalesOrder2", docs[0].Id)
}

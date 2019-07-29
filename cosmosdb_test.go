package gocosmosdb

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

// for advanced debugging
// client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", Debug: true, Verbose: true}, log)
func TestNew(t *testing.T) {
	assert := assert.New(t)
	client := New("url", Config{MasterKey: "config"}, log)
	assert.IsType(client, &CosmosDB{}, "Should return CosmosDB object")
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
	s := ServerFactory(resp, 200)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg==", Debug: true, Verbose: true}, log)
	//client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	db, err := client.ReadDatabase("dbs/qicAAA==")
	spew.Dump(db)
	assert.Nil(err)
	assert.Nil(db)
	assert.Equal("iot2", db.Id)
}

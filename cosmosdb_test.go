package gocosmosdb

import (
	"testing"

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
	s := ServerFactory(`{"id": "db-id"}`, 200)
	defer s.Close()
	client := New(s.URL, Config{MasterKey: "YXJpZWwNCg=="}, log)
	db, err := client.ReadDatabase("dbs/b7NTAS==/")
	assert.Nil(err)
	assert.NotNil(db)
	assert.Equal("db-id", db.Id)
}

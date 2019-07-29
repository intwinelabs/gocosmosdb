# gocosmosdb -  Go client for Azure CosmosDB SQL API
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/intwinelabs/gocosmosdb)
[![Build Status](https://travis-ci.org/intwinelabs/gocosmosdb.svg?branch=master)](https://travis-ci.org/intwinelabs/gocosmosdb)
[![Coverage Status](https://coveralls.io/repos/github/intwinelabs/gocosmosdb/badge.svg?branch=master)](https://coveralls.io/github/intwinelabs/gocosmosdb?branch=master)

### Key Features
- Client Connection Pooling
- Retry With Backoff
- TTL for documents
- Advanced Debugging

### Get Started

#### Installation
```bash
$ go get github.com/intwinelabs/gocosomsdb
```

#### Example
```go
import (
	"io/ioutil"
	"log"

	"github.com/intwinelabs/gocosmosdb"
	"github.com/intwinelabs/logger"
)

func main() {
	log := logger.Init("CosmosGoApp", false, true, ioutil.Discard)
	client := gocosmosdb.New("connection-url", gocosmosdb.Config{MasterKey: "master-key"}, log)
	
	// create a database
	db, err := client.CreateDatabase(`{ "id": "intwineLabs" }`)
	if err != nil {
		log.Fatal(err)
	}

	// create a collection
	coll, err := client.CreateCollection(db.Self, `{"id": "users"}`)
	if err != nil {
		log.Fatal(err)
	}

	// user struct
	type User struct {
		gocosmosdb.Document
		// To set documents TTL
		gocosmosdb.Expirable
		// Your external fields
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
	}

	// user to store
	var user User
	// Note: If the `Id` is missing(or empty) in the payload it will generate 
	// random document id(i.e: uuid4)
	user.Id = "uuid"
	user.Name = "Bad MF"
	user.Email = "badmf@intwine.io"
	// This tells CosmosDB to expire the doc in 24 hours
	user.SetTTL(24 * time.Hour)
	
	// create the document
	err := client.CreateDocument(col.Self, &doc)
	if err != nil {
		log.Fatal(err)	
	}

	// query to documents
	var users []User
	err = client.QueryDocuments(coll.Self, "SELECT * FROM ROOT r", &users)
	if err != nil {
		log.Fatal(err)	
	}
	for _, user := range users {
		fmt.Print("Name:", user.Name, "Email:", user.Email)
	}
}
```

### Credits
This library is derived from:
- https://github.com/nerdylikeme/go-documentdb
- https://github.com/a8m/documentdb

### License

Copyright (c) 2018 Intwine Labs, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
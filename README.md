# gocosmosdb -  Go client for Azure CosmosDB SQL API
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/intwinelabs/gocosmosdb)
[![Build Status](https://travis-ci.org/intwinelabs/gocosmosdb.svg?branch=master)](https://travis-ci.org/intwinelabs/gocosmosdb)
[![Coverage Status](https://coveralls.io/repos/github/intwinelabs/gocosmosdb/badge.svg?branch=master)](https://coveralls.io/github/intwinelabs/gocosmosdb?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/intwinelabs/gremgoser)](https://goreportcard.com/report/github.com/intwinelabs/gremgoser)

### Key Features
- Client Connection Pooling
- Retry With Backoff
- TTL for documents
- Advanced Debugging

### Get Started

#### Installation
```bash
$ go get github.com/intwinelabs/gocosmosdb
```

#### Example
```go
package main

import (
	"context"
	"io/ioutil"
	"time"

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
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
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
	_, err = client.CreateDocument(coll.Self, &user)
	if err != nil {
		log.Fatal(err)
	}

	// query to documents
	var users []User
	_, err = client.QueryDocuments(coll.Self, "SELECT * FROM ROOT r", &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Infof("Name:%s, Email: %s", user.Name, user.Email)
	}

	// run stored procedure with context timeout
	ctx, _ := context.WithTimeout(context.Background(), 250*time.Millisecond)
	docs := []User{}
	_, err = client.ExecuteStoredProcedure(coll.Self+"sprocs/Sl8fALN4sw4CAAAAAAAAgA==", []string{"param1"}, &docs, gocosmosdb.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}
}
```

### Azure Cosmos DB SQL REST API Reference
- https://docs.microsoft.com/en-us/rest/api/cosmos-db/

### Contributing
- **Reporting Issues** - When reporting issues on GitHub please include your host OS (Ubuntu 16.04, Fedora 19, etc) `sudo lsb_release -a`, the output of `uname -a`, `go version`. Please include the steps required to reproduce the problem. This info will help us review and fix your issue faster.
- **We welcome your pull requests** - We are always thrilled to receive pull requests, and do our best to process them as fast as possible. 
	- Not sure if that typo is worth a pull request? Do it! We will appreciate it.
    - If your pull request is not accepted on the first try, don't be discouraged! We will do our best to give you feedback on what to improve.
    - We're trying very hard to keep gocosmosdb lean and focused. We don't want it to do everything for everybody. This means that we might decide against incorporating a new feature. However, we encourage you to fork our repo and implement it on top of gocosmosdb.
	- Any changes or improvements should be documented as a GitHub issue before we add it to the project and anybody starts working on it.
- **Please check for existing issues first** - If it does add a quick "+1". This will help prioritize the most common problems and requests.

### Credits
This library is derived from:
- https://github.com/nerdylikeme/go-documentdb
- https://github.com/a8m/documentdb

### License

Copyright (c) 2018 Intwine Labs, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

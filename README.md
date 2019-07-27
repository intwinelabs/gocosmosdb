# gocosmosdb
> Go client for Azure CosmosDB SQL API

### Key Features
- Client HTTP Connection Pooling 
- Retry With Backoff
- TTL for documents
- Advanced Debugging

### Credits
This library is derived from:
- https://github.com/nerdylikeme/go-documentdb
- https://github.com/a8m/documentdb

### Get Started

#### Installation
```bash
$ go get github.com/intwinelabs/gocosomsdb
```

#### Add to your project
```go
import (
	"github.com/intwinelabs/gocosmosdb"
)

func main() {
	client := gocosmosdb.New("connection-url", gocosmosdb.Config{"master-key"})
	// Start using gocosmosdb
	dbs, err := client.ReadDatabases()
	if err != nill {
	  log.Fatal(err)
	}
	fmt.Println(dbs)
}
```

#### QueryDatabases
```go
dbs, err := client.QueryDatabases("SELECT * FROM ROOT r")
if err != nil {
	log.Fatal(err)	
}
for _, db := range dbs {
	fmt.Println("DB Name:", db.Id)
}
```
[TOC](#Table_of_contents)


#### CreateDatabase
```go
db, err := client.CreateDatabase(`{ "id": "test" }`)
if err != nil {
	log.Fatal(err)	
}
fmt.Println(db)
```
```go	
// or ...
var db gocosmosdb.Database
db.Id = "test"
db, err = client.CreateDatabase(&db)
```
[TOC](#Table_of_contents)

#### ReplaceDatabase
```go
db, err := client.ReplaceDatabase("self_link | constructed_uri", `{ "id": "test" }`)
if err != nil {
	log.Fatal(err)	
}
fmt.Println(db)
```
```go	
// or ...
var db gocosmosdb.Database
db, err = client.ReplaceDatabase("self_link | constructed_uri", &db)
```
[TOC](#Table_of_contents)

#### DeleteDatabase
```go
err := client.DeleteDatabase("self_link | constructed_uri")
if err != nil {
	log.Fatal(err)	
}
```
[TOC](#Table_of_contents)

### Collections

#### QueryCollections
```go
colls, err := client.QueryCollections("db_self_link | constructed_uri", "SELECT * FROM ROOT r")
if err != nil {
	log.Fatal(err)	
}
for _, coll := range colls {
	fmt.Println("Collection Name:", coll.Id)
}
```
[TOC](#Table_of_contents)

#### ReadCollections
```go
colls, err := client.ReadCollections("db_self_link | constructed_uri")
if err != nil {
	log.Fatal(err)	
}
for _, coll := range colls {
	fmt.Println("Collection Name:", coll.Id)
}
```
[TOC](#Table_of_contents)

#### CreateCollection
```go
coll, err := client.CreateCollection("db_self_link | constructed_uri", `{"id": "my_test"}`)
if err != nil {
	log.Fatal(err)	
}
fmt.Println("Collection Name:", coll.Id)
```
```go	
// or ...
var coll gocosmosdb.Collection
coll.Id = "test"
coll, err = client.CreateCollection("db_self_link | constructed_uri", &coll)
```
[TOC](#Table_of_contents)

#### DeleteCollection
```go
err := client.DeleteCollection("self_link | constructed_uri")
if err != nil {
	log.Fatal(err)	
}
```
[TOC](#Table_of_contents)


#### QueryDocuments
```go
type User struct {
	gocosmosdb.Document
	// Your external fields
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
}

func main() {
	var users []User
	err = client.QueryDocuments("coll_self_link | constructed_uri", "SELECT * FROM ROOT r", &users)
	if err != nil {
		log.Fatal(err)	
	}
	for _, user := range users {
		fmt.Print("Name:", user.Name, "Email:", user.Email)
	}
}
```
[TOC](#Table_of_contents)

#### ReDocuments
```go
type User struct {
	gocosmosdb.Document
	// Your external fields
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
}

func main() {
	var users []User
	err = client.ReadDocuments("coll_self_link | constructed_uri", &users)
	if err != nil {
		log.Fatal(err)	
	}
	for _, user := range users {
		fmt.Print("Name:", user.Name, "Email:", user.Email)
	}
}
```
[TOC](#Table_of_contents)

#### CreateDocument
```go
type User struct {
	gocosmosdb.Document
	// To set docuemnts TTL
	gocosmosdb.Expirable
	// Your external fields
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
}

func main() {
	var user User
	// Note: If the `id` is missing(or empty) in the payload it will generate 
	// random document id(i.e: uuid4)
	user.Id = "uuid"
	user.Name = "Bad MF"
	user.Email = "badmf@intwine.io"
	// This tells CosmosDB to expire the doc in 24 hours
	user.SetTTL(24 * time.Hour)
	err := client.CreateDocument("coll_self_link | constructed_uri", &doc)
	if err != nil {
		log.Fatal(err)	
	}
	fmt.Print("Name:", user.Name, "Email:", user.Email)
}
```
[TOC](#Table_of_contents)

#### ReplaceDocument
```go
type User struct {
	gocosmosdb.Document
	// Your external fields
	IsAdmin bool   `json:"isAdmin,omitempty"`
}

func main() {
	var user User
	user.Id = "uuid"
	user.IsAdmin = false
	err := client.ReplaceDocument("doc_self_link | constructed_uri", &user)
	if err != nil {
		log.Fatal(err)	
	}
	fmt.Print("Is Admin:", user.IsAdmin)
}
```
[TOC](#Table_of_contents)

#### DeleteDocument
```go
err := client.DeleteDocument("doc_self_link | constructed_uri")
if err != nil {
	log.Fatal(err)	
}
```
[TOC](#Table_of_contents)



#### ExecuteStoredProcedure
```go
var docs []Document
err := client.ExecuteStoredProcedure("sporc_self", [...]interface{}{p1, p2}, &docs)
if err != nil {
	log.Fatal(err)
}
```


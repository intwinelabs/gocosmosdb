# gocosmosdb
> Go client driver for Azure CosmosDB 

### Note
This library is derived from `github.com/nerdylikeme/go-documentdb`(github.com/nerdylikeme/go-documentdb)

TODO: Add documentation on stored procedurs and user defined functions

## Table of contents:
- [Get Started](#get-started)
- [Examples](#examples)
- [Databases](#databases)
  - [Get](#readdatabase)
  - [Query](#querydatabases)
  - [List](#readdatabases)
  - [Create](#createdatabase)
  - [Replace](#replacedatabase)
  - [Delete](#deletedatabase)
- [Collections](#collections)
  - [Get](#readcollection)
  - [Query](#querycollections)
  - [List](#readcollection)
  - [Create](#createcollection)
  - [Delete](#deletecollection)
- [Documents](#documents)
  - [Get](#readdocument)
  - [Query](#querydocuments)
  - [List](#readdocuments)
  - [Create](#createdocument)
  - [Replace](#replacedocument)
  - [Delete](#deletedocument)
- [StoredProcedures](#storedprocedures)
  - [Get](#readstoredprocedure)
  - [Query](#querystoredprocedures)
  - [List](#readstoredprocedures)
  - [Create](#createstoredprocedure)
  - [Replace](#replacestoredprocedure)
  - [Delete](#deletestoredprocedure)
  - [Execute](#executestoredprocedure)
- [UserDefinedFunctions](#userdefinedfunctions)
  - [Get](#readuserdefinedfunction)
  - [Query](#queryuserdefinedfunctions)
  - [List](#readuserdefinedfunctions)
  - [Create](#createuserdefinedfunction)
  - [Replace](#replaceuserdefinedfunction)
  - [Delete](#deleteuserdefinedfunction)

### Get Started

#### Installation
```bash
$ go get github.com/intwineapp/gocosomsdb
```
[TOC](#Table_of_contents)

#### Add to your project
```go
import (
	"github.com/intwineapp/gocosmosdb"
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
[TOC](#Table_of_contents)

### Databases

#### ReadDatabase
```go
db, err := client.ReadDatabase("self_link")
if err != nil {
	log.Fatal(err)	
}
fmt.Println(db.Self, db.Id)
```
[TOC](#Table_of_contents)

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

#### ReadDatabases
```go
dbs, err := client.ReadDatabases()
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
db, err := client.ReplaceDatabase("self_link", `{ "id": "test" }`)
if err != nil {
	log.Fatal(err)	
}
fmt.Println(db)
```
```go	
// or ...
var db gocosmosdb.Database
db, err = client.ReplaceDatabase("self_link", &db)
```
[TOC](#Table_of_contents)

#### DeleteDatabase
```go
err := client.DeleteDatabase("self_link")
if err != nil {
	log.Fatal(err)	
}
```
[TOC](#Table_of_contents)

### Collections

#### ReadCollection
```go
coll, err := client.ReadCollection("self_link")
if err != nil {
	log.Fatal(err)	
}
fmt.Println(coll.Self, coll.Id)
```
[TOC](#Table_of_contents)

#### QueryCollections
```go
colls, err := client.QueryCollections("db_self_link", "SELECT * FROM ROOT r")
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
colls, err := client.ReadCollections("db_self_link")
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
coll, err := client.CreateCollection("db_self_link", `{"id": "my_test"}`)
if err != nil {
	log.Fatal(err)	
}
fmt.Println("Collection Name:", coll.Id)
```
```go	
// or ...
var coll gocosmosdb.Collection
coll.Id = "test"
coll, err = client.CreateCollection("db_self_link", &coll)
```
[TOC](#Table_of_contents)

#### DeleteCollection
```go
err := client.DeleteCollection("self_link")
if err != nil {
	log.Fatal(err)	
}
```
[TOC](#Table_of_contents)

### Documents

#### ReadDocument
```go
type Document struct {
	gocosmosdb.Document
	// Your external fields
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
}

func main() {
	var doc Document
	err = client.ReadDocument("self_link", &doc)
	if err != nil {
		log.Fatal(err)	
	}
	fmt.Println("Document Name:", doc.Name)
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
	err = client.QueryDocuments("coll_self_link", "SELECT * FROM ROOT r", &users)
	if err != nil {
		log.Fatal(err)	
	}
	for _, user := range users {
		fmt.Print("Name:", user.Name, "Email:", user.Email)
	}
}
```
[TOC](#Table_of_contents)

#### ReadDocuments
```go
type User struct {
	gocosmosdb.Document
	// Your external fields
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
}

func main() {
	var users []User
	err = client.ReadDocuments("coll_self_link", &users)
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
	err := client.CreateDocument("coll_self_link", &doc)
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
	err := client.ReplaceDocument("doc_self_link", &user)
	if err != nil {
		log.Fatal(err)	
	}
	fmt.Print("Is Admin:", user.IsAdmin)
}
```
[TOC](#Table_of_contents)

#### DeleteDocument
```go
err := client.DeleteDocument("doc_self_link")
if err != nil {
	log.Fatal(err)	
}
```
[TOC](#Table_of_contents)

### StoredProcedures

#### ReadStoredProcedure
```go
client.ReadStoredProcedure()
```
[TOC](#Table_of_contents)

#### QueryStoredProcedure
```go
client.QueryStoredProcedure()
```
[TOC](#Table_of_contents)

#### CreateStoredProcedure
```go
client.ReadStoredProcedure()
```
[TOC](#Table_of_contents)

#### ReplaceStoredProcedure
```go
client.ReplaceStoredProcedure()
```
[TOC](#Table_of_contents)

#### DeleteStoredProcedure
```go
client.DeleteStoredProcedure()
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
[TOC](#Table_of_contents)

### UserDefinedFunctions

#### ReadUserDefinedFunction
```go
client.ReadUserDefinedFunction()
```
[TOC](#Table_of_contents)

#### QueryUserDefinedFunction
```go
client.QueryUserDefinedFunction()
```
[TOC](#Table_of_contents)

#### CreateUserDefinedFunction
```go
client.ReadUserDefinedFunction()
```
[TOC](#Table_of_contents)

#### ReplaceUserDefinedFunction
```go
client.ReplaceUserDefinedFunction()
```
[TOC](#Table_of_contents)

#### DeleteUserDefinedFunction
```go
client.DeleteUserDefinedFunction()
```
[TOC](#Table_of_contents)

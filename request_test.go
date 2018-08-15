package gocosmosdb

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceRequest(t *testing.T) {
	assert := assert.New(t)
	req := ResourceRequest("/dbs/b5NCAA==/", &http.Request{})
	assert.Equal(req.rType, "dbs")
	assert.Equal(req.rId, "b5NCAA==")
}

func TestDefaultHeaders(t *testing.T) {
	r, _ := http.NewRequest("GET", "link", &bytes.Buffer{})
	req := ResourceRequest("/dbs/b5NCAA==/", r)
	_ = req.DefaultHeaders("YXJpZWwNCg==")

	assert := assert.New(t)
	assert.NotEqual(req.Header.Get(HEADER_AUTH), "")
	assert.NotEqual(req.Header.Get(HEADER_XDATE), "")
	assert.NotEqual(req.Header.Get(HEADER_VER), "")
}

// Get link and return resource Id and Type
// /dbs	Feed of databases under a database account - 1 - 3
// /dbs/{dbName}	Database with an id matching the value {dbName} - 2 - 4
// /dbs/{dbName}/colls/	Feed of collections under a database - 3 - 5
// /dbs/{dbName}/colls/{collName}	Collection with an id matching the value {collName} - 4 - 6
// /dbs/{dbName}/colls/{collName}/docs	Feed of documents under a collection - 5 - 7
// /dbs/{dbName}/colls/{collName}/docs/{docId}	Document with an id matching the value {doc} - 6 - 8
// /dbs/{dbName}/users/	Feed of users under a database - 3 - 5
// /dbs/{dbName}/users/{userId}	User with an id matching the value {user} - 4 -6
// /dbs/{dbName}/users/{userId}/permissions	Feed of permissions under a user - 5 -7
// /dbs/{dbName}/users/{userId}/permissions/{permissionId}	Permission with an id matching the value {permission} - 6 - 8
// (e.g: "/dbs/b5NCAA==/" ==> "b5NCAA==", "b5NCAA==", "dbs")
// (e.g: "/dbs/mydb/colls/mydb/docs/mydoc" ==> "b5NCAA==", "docs")
func TestParseLink(t *testing.T) {
	assert := assert.New(t)

	// /dbs	Feed of databases under a database account - 1 - 3
	link := "/dbs"
	rLink, rId, rType := parse(link)
	assert.Equal("", rLink)
	assert.Equal("", rId)
	assert.Equal("dbs", rType)

	// /dbs/{dbName}	Database with an id matching the value {dbName} - 2 - 4
	link = "/dbs/b5NCAA=="
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAA==", rLink)
	assert.Equal("b5NCAA==", rId)
	assert.Equal("dbs", rType)

	// /dbs/{dbName}	Database with an id matching the value {dbName} - 2 - 4
	link = "/dbs/mydb"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb", rLink)
	assert.Equal("mydb", rId)
	assert.Equal("dbs", rType)

	// /dbs/{dbName}/colls/	Feed of collections under a database - 3 - 5
	link = "/dbs/b5NCAA==/colls/"
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAA==", rLink)
	assert.Equal("b5NCAA==", rId)
	assert.Equal("colls", rType)

	// /dbs/{dbName}/colls/	Feed of collections under a database - 3 - 5
	link = "/dbs/mydb/colls/"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/colls", rLink)
	assert.Equal("", rId)
	assert.Equal("colls", rType)

	// /dbs/{dbName}/colls/{collName}	Collection with an id matching the value {collName} - 4 - 6
	link = "/dbs/b5NCAA==/colls/b5NCAB=="
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAB==", rLink)
	assert.Equal("b5NCAB==", rId)
	assert.Equal("colls", rType)

	// /dbs/{dbName}/colls/{collName}	Collection with an id matching the value {collName} - 4 - 6
	link = "/dbs/mydb/colls/mycoll"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/colls/mycoll", rLink)
	assert.Equal("mycoll", rId)
	assert.Equal("colls", rType)

	// /dbs/{dbName}/colls/{collName}/docs	Feed of documents under a collection - 5 - 7
	link = "/dbs/b5NCAA==/colls/b5NCAB==/docs"
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAB==", rLink)
	assert.Equal("b5NCAB==", rId)
	assert.Equal("docs", rType)

	// /dbs/{dbName}/colls/{collName}/docs	Feed of documents under a collection - 5 - 7
	link = "/dbs/mydb/colls/mycoll/docs"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/colls/mycoll/docs", rLink)
	assert.Equal("", rId)
	assert.Equal("docs", rType)

	// /dbs/{dbName}/colls/{collName}/docs/{docId}	Document with an id matching the value {doc} - 6 - 8
	link = "/dbs/b5NCAA==/colls/b5NCAB==/docs/b5NCAC=="
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAC==", rLink)
	assert.Equal("b5NCAC==", rId)
	assert.Equal("docs", rType)

	// /dbs/{dbName}/colls/{collName}/docs/{docId}	Document with an id matching the value {doc} - 6 - 8
	link = "/dbs/mydb/colls/mycoll/docs/mydoc"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/colls/mycoll/docs/mydoc", rLink)
	assert.Equal("mydoc", rId)
	assert.Equal("docs", rType)

	// /dbs/{dbName}/users/	Feed of users under a database - 3 - 5
	link = "/dbs/b5NCAA==/users/"
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAA==", rLink)
	assert.Equal("b5NCAA==", rId)
	assert.Equal("users", rType)

	// /dbs/{dbName}/users/	Feed of users under a database - 3 - 5
	link = "/dbs/mydb/users/"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/users", rLink)
	assert.Equal("", rId)
	assert.Equal("users", rType)

	// /dbs/{dbName}/users/{userId}	User with an id matching the value {user} - 4 -6
	link = "/dbs/b5NCAA==/users/b5NCAB=="
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAB==", rLink)
	assert.Equal("b5NCAB==", rId)
	assert.Equal("users", rType)

	// /dbs/{dbName}/users/{userId}	User with an id matching the value {user} - 4 -6
	link = "/dbs/mydb/users/mycoll"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/users/mycoll", rLink)
	assert.Equal("mycoll", rId)
	assert.Equal("users", rType)

	// /dbs/{dbName}/users/{userId}/permissions	Feed of permissions under a user - 5 -7
	link = "/dbs/b5NCAA==/users/b5NCAB==/permissions"
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAB==", rLink)
	assert.Equal("b5NCAB==", rId)
	assert.Equal("permissions", rType)

	// /dbs/{dbName}/users/{userId}/permissions	Feed of permissions under a user - 5 -7
	link = "/dbs/mydb/users/mycoll/permissions"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/users/mycoll/permissions", rLink)
	assert.Equal("", rId)
	assert.Equal("permissions", rType)

	// /dbs/{dbName}/users/{userId}/permissions/{permissionId}	Permission with an id matching the value {permission} - 6 - 8
	link = "/dbs/b5NCAA==/users/b5NCAB==/permissions/b5NCAC=="
	rLink, rId, rType = parse(link)
	assert.Equal("b5NCAC==", rLink)
	assert.Equal("b5NCAC==", rId)
	assert.Equal("permissions", rType)

	// /dbs/{dbName}/users/{userId}/permissions/{permissionId}	Permission with an id matching the value {permission} - 6 - 8
	link = "/dbs/mydb/users/mycoll/permissions/mydoc"
	rLink, rId, rType = parse(link)
	assert.Equal("dbs/mydb/users/mycoll/permissions/mydoc", rLink)
	assert.Equal("mydoc", rId)
	assert.Equal("permissions", rType)
}

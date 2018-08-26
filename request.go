package gocosmosdb

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HEADER_XDATE             = "X-Ms-Date"
	HEADER_AUTH              = "Authorization"
	HEADER_VER               = "X-Ms-Version"
	HEADER_CONTYPE           = "Content-Type"
	HEADER_CONLEN            = "Content-Length"
	HEADER_IS_QUERY          = "X-Ms-Documentdb-Isquery"
	HEADER_POP_QUERY_METRICS = "X-Ms-Documentdb-Populatequerymetrics"
	HEADER_QUERY_METRICS     = "X-Ms-Documentdb-Query-Metrics"
	HEADER_REQ_CHARGE        = "X-Ms-Request-Charge"
)

// Request Error
type RequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Implement Error function
func (e RequestError) Error() string {
	return fmt.Sprintf("%v, %v", e.Code, e.Message)
}

// Resource Request
type Request struct {
	rLink string
	rId   string
	rType string
	*http.Request
}

// Return new resource request with type and id
func ResourceRequest(link string, req *http.Request) *Request {
	rLink, rId, rType := parse(link)
	return &Request{rLink, rId, rType, req}
}

// Add 3 default headers to *Request
// "x-ms-date", "x-ms-version", "authorization"
func (req *Request) DefaultHeaders(mKey string) (err error) {
	req.Header.Add(HEADER_XDATE, time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	req.Header.Add(HEADER_VER, "2017-02-22")

	// Auth
	parts := req.Method + "\n" +
		req.rType + "\n" +
		req.rLink + "\n" +
		req.Header.Get(HEADER_XDATE) + "\n" +
		req.Header.Get("Date") + "\n"

	partsLower := strings.ToLower(parts)

	sign, err := authorize(partsLower, mKey)
	if err != nil {
		return err
	}

	masterToken := "master"
	tokenVersion := "1.0"
	req.Header.Add(HEADER_AUTH, url.QueryEscape("type="+masterToken+"&ver="+tokenVersion+"&sig="+sign))
	return
}

// Add headers for query request
func (req *Request) QueryHeaders(len int) {
	req.Header.Add(HEADER_CONTYPE, "application/query+json")
	req.Header.Add(HEADER_IS_QUERY, "true")
	req.Header.Add(HEADER_CONLEN, string(len))
}

// Add headers for query metrics request
func (req *Request) QueryMetricsHeaders() {
	req.Header.Add(HEADER_POP_QUERY_METRICS, "true")
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
func parse(link string) (rLink, rId, rType string) {
	if strings.HasPrefix(link, "/") == false {
		link = "/" + link
	}
	if strings.HasSuffix(link, "/") == false {
		link = link + "/"
	}

	parts := strings.Split(link, "/")
	l := len(parts)

	//spew.Dump(parts)
	if strings.Index(parts[2], "==") > -1 { // use this logic if it's a _self link
		if l%2 == 0 {
			rLink = parts[l-2]
			rId = parts[l-2]
			rType = parts[l-3]
		} else {
			rLink = parts[l-3]
			rId = parts[l-3]
			rType = parts[l-2]
		}
	} else { // use this logic if it's a constructed uri using ids
		if l == 3 && parts[1] == "dbs" {
			rLink = ""
			rId = ""
			rType = parts[1]
		} else if l == 4 && parts[1] == "dbs" {
			rLink = parts[1] + "/" + parts[2]
			rId = parts[2]
			rType = parts[1]
		} else if l == 5 && parts[1] == "dbs" && (parts[3] == "colls" || parts[3] == "users") {
			rLink = parts[1] + "/" + parts[2] + "/" + parts[3]
			rId = ""
			rType = parts[3]
		} else if l == 6 && parts[1] == "dbs" && (parts[3] == "colls" || parts[3] == "users") {
			rLink = parts[1] + "/" + parts[2] + "/" + parts[3] + "/" + parts[4]
			rId = parts[4]
			rType = parts[3]
		} else if l == 7 && parts[1] == "dbs" {
			rLink = parts[1] + "/" + parts[2] + "/" + parts[3] + "/" + parts[4] + "/" + parts[5]
			rId = ""
			rType = parts[5]
		} else if l == 8 && parts[1] == "dbs" {
			rLink = parts[1] + "/" + parts[2] + "/" + parts[3] + "/" + parts[4] + "/" + parts[5] + "/" + parts[6]
			rId = parts[6]
			rType = parts[5]
		}
	}

	return
}

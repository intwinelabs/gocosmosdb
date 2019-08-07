package gocosmosdb

import "errors"

// NewPagableQuery - Creates a pagable query that populates the passed docs interface
func (c *CosmosDB) NewPagableQuery(coll string, query *QueryWithParameters, limit int, docs interface{}, opts ...CallOption) *PagableQuery {
	return &PagableQuery{
		client: c,
		coll:   coll,
		query:  query,
		limit:  Limit(limit),
		offset: 0,
		docs:   docs,
		opts:   opts,
	}
}

func (q *PagableQuery) doQuery(coll string, query *QueryWithParameters, docs interface{}, opts ...CallOption) (*Response, error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if query != nil {
		return q.client.client.queryWithParameters(coll+"docs/", query, &data, opts...)
	}
	return nil, errors.New("QueryWithParameters cannot be nil")
}

// Next - marshals the next page of docs into the passed interface
func (q *PagableQuery) Next() error {
	if q.offset > 0 {
		opts := append(q.opts, q.limit)
		opts = append(opts, q.continuation)
		opts = append(opts, q.sessionToken)
		resp, err := q.doQuery(q.coll, q.query, q.docs, opts...)
		if err != nil {
			return err
		}
		q.offset = q.offset + 1
		q.continuation = Continuation(resp.Continuation())
	}
	if q.offset == 0 {
		opts := append(q.opts, q.limit)
		resp, err := q.doQuery(q.coll, q.query, q.docs, opts...)
		if err != nil {
			return err
		}
		q.offset = q.offset + 1
		q.continuation = Continuation(resp.Continuation())
		q.sessionToken = SessionToken(resp.SessionToken())
	}
	return nil
}

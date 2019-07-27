package gocosmosdb

import "errors"

// NewPagableQuery
func (c *CosmosDB) NewPagableQuery(coll string, query *QueryWithParameters, limit int, docs interface{}) *PagableQuery {
	return &PagableQuery{
		client: c,
		coll:   coll,
		query:  query,
		limit:  Limit(limit),
		offset: 0,
		docs:   docs,
	}
}

func (q *PagableQuery) doQuery(coll string, query *QueryWithParameters, docs interface{}, opts ...CallOption) (*Response, error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if query != nil {
		return q.client.client.QueryWithParameters(coll+"docs/", query, &data, opts...)
	}
	return nil, errors.New("QueryWithParameters cannot be nil")
}

// Next
func (q *PagableQuery) Next() error {
	if q.offset > 0 {
		resp, err := q.doQuery(q.coll, q.query, q.docs, q.limit, q.continuation, q.sessionToken)
		if err != nil {
			return err
		}
		q.offset = q.offset + 1
		q.continuation = Continuation(resp.Continuation())
	}
	if q.offset == 0 {
		resp, err := q.doQuery(q.coll, q.query, q.docs, q.limit)
		if err != nil {
			return err
		}
		q.offset = q.offset + 1
		q.continuation = Continuation(resp.Continuation())
		q.sessionToken = SessionToken(resp.SessionToken())
	}
	return nil
}
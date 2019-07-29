package gocosmosdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpectStatusCodes(t *testing.T) {
	assert := assert.New(t)

	st1 := expectStatusCode(200)
	assert.Equal(true, st1(200))

	st2 := expectStatusCodeXX(400)
	assert.Equal(true, st2(499))

}

package gocosmosdb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert := assert.New(t)

	uuid := genId()
	assert.Equal(36, len(uuid))

	exp := Expirable{}
	exp.SetTTL(100 * time.Second)
	assert.Equal(int64(100), exp.TTL)

	b, err := stringify([]byte("foo"))
	assert.Nil(err)
	assert.Equal([]byte("foo"), b)
}

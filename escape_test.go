package gocosmosdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeSQL(t *testing.T) {
	assert := assert.New(t)

	sql := escapeSQL("SELECT * FROM root \n\r\x00\x1a\" \\ r")
	assert.Equal("SELECT * FROM root \\n\\r\\0\\Z\\\" \\\\ r", sql)
}

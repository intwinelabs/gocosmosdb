package gocosmosdb

import "github.com/google/uuid"

func genId() string {
	return uuid.New().String()
}

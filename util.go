package gocosmosdb

import (
	"math"
	"time"

	"github.com/google/uuid"
)

func genId() string {
	return uuid.New().String()
}

// SetTTL takes a duration and sets the field value
func (exp *Expirable) SetTTL(dur time.Duration) {
	exp.TTL = int64(math.Round(dur.Seconds()))
}

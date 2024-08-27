package model

import (
	"time"
	"xrf197ilz35aq0/random"
)

type RungTrail struct {
	id       int64
	score    int
	created  time.Time
	tradeIds map[string]bool // a set of Ids
}

func NewRungTrail() *RungTrail {
	id := random.PositiveInt64()
	return &RungTrail{
		id:       id,
		score:    0,
		created:  time.Now(),
		tradeIds: make(map[string]bool),
	}
}

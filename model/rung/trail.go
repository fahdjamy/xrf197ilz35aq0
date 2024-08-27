package rung

import (
	"time"
	"xrf197ilz35aq0/random"
)

type TrailMetaData struct {
	rating int
}

type Trail struct {
	id       int64
	score    int
	created  time.Time
	tradeIds map[string]TrailMetaData // a set metadata trade trails
}

func (rt *Trail) UpdateMetaData(tradeId string, customerRating int) {
	trailMetaData, ok := rt.tradeIds[tradeId]
	if ok {
		trailMetaData.rating = customerRating
	} else {
		rt.tradeIds[tradeId] = TrailMetaData{
			rating: customerRating,
		}
	}
}

func (rt *Trail) Score() int {
	return rt.score
}

func NewRungTrail() *Trail {
	id := random.PositiveInt64()
	return &Trail{
		id:       id,
		score:    0,
		created:  time.Now(),
		tradeIds: make(map[string]TrailMetaData),
	}
}

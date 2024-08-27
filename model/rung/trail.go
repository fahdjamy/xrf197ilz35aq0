package rung

import (
	"time"
	"xrf197ilz35aq0/random"
)

type RungTrailMeta struct {
	rating int
}

type RungTrail struct {
	id       int64
	score    int
	created  time.Time
	tradeIds map[string]RungTrailMeta // a set of Ids
}

func (rt *RungTrail) UpdateMetaData(tradeId string, customerRating int) {
	trailMetaData, ok := rt.tradeIds[tradeId]
	if ok {
		trailMetaData.rating = customerRating
	} else {
		rt.tradeIds[tradeId] = RungTrailMeta{
			rating: customerRating,
		}
	}
}

func (rt *RungTrail) Score() int {
	return rt.score
}

func NewRungTrail() *RungTrail {
	id := random.PositiveInt64()
	return &RungTrail{
		id:       id,
		score:    0,
		created:  time.Now(),
		tradeIds: make(map[string]RungTrailMeta),
	}
}

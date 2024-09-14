package rung

import (
	"fmt"
	"time"
	"xrf197ilz35aq0/core/model"
	"xrf197ilz35aq0/internal/random"
)

type TrailMetaData struct {
	Rating  int
	Comment string
}

func (t *TrailMetaData) String() string {
	return fmt.Sprintf("TrailMetaData{Rating: %d, Comment: %s}", t.Rating, t.Comment)
}

func (t *TrailMetaData) IsEqual(other TrailMetaData) bool {
	return t.Comment == other.Comment && t.Rating == other.Rating
}

const (
	scoreMin = 0
	scoreMax = 10
)

type Trail struct {
	id       int64
	score    int
	created  time.Time
	updated  time.Time
	tradeIds map[string]TrailMetaData // a set metadata trade trails
}

func (rt *Trail) UpdateMetaData(tradeId string, metadata TrailMetaData) (bool, error) {
	if metadata.Rating < scoreMin || metadata.Rating > scoreMax {
		return false, model.InvalidRequest{Message: "Trail Rating out of bounds"}
	}

	trailMetaData, ok := rt.tradeIds[tradeId]
	isUpdated := false
	// doesn't exist
	if !ok {
		rt.tradeIds[tradeId] = metadata
		isUpdated = true
	} else if !trailMetaData.IsEqual(metadata) {
		trailMetaData.Rating = metadata.Rating
		trailMetaData.Comment = metadata.Comment
		isUpdated = true
	}

	if isUpdated {
		rt.updated = time.Now()
		rt.calculateScore()
	}

	return isUpdated, nil
}

func (rt *Trail) Score() int {
	return rt.score
}

func (rt *Trail) UpdatedAt() time.Time {
	return rt.updated
}

func (rt *Trail) calculateScore() {
	newScore := 0
	if len(rt.tradeIds) > 0 {
		for _, trailMetaData := range rt.tradeIds {
			newScore += trailMetaData.Rating
		}
	}
	// average of the score
	rt.score = newScore / len(rt.tradeIds)
}

func NewRungTrail() *Trail {
	id := random.PositiveInt64()
	now := time.Now()
	return &Trail{
		id:       id,
		score:    0,
		created:  now,
		updated:  now,
		tradeIds: make(map[string]TrailMetaData),
	}
}

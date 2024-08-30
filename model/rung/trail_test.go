package rung_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/model/rung"
)

var validMetadataTrail = rung.TrailMetaData{
	Rating:  1,
	Comment: "This is the comment",
}

var maxRatingMetadataTrail = rung.TrailMetaData{
	Rating:  100,
	Comment: "This is the comment",
}

var tradeId = "someId"

func TestNewRungTrail(t *testing.T) {
	t.Run("should return a new RungTrail struct", func(t *testing.T) {
		trail := rung.NewRungTrail()
		assertExpectedTrailScore(t, *trail, 0)
	})
}

func TestRungTrailUpdateMetaData(t *testing.T) {
	t.Run("should update the metadata", func(t *testing.T) {
		trail := rung.NewRungTrail()
		updated, err := trail.UpdateMetaData(tradeId, validMetadataTrail)
		xrf197ilz35aq0.AssertNoError(t, err)
		assert.True(t, updated)
	})

	t.Run("should return an error if rating is out of bound", func(t *testing.T) {
		trail := rung.NewRungTrail()
		_, err := trail.UpdateMetaData(tradeId, maxRatingMetadataTrail)
		xrf197ilz35aq0.AssertError(t, err)
	})

	t.Run("should return not update if trail exists", func(t *testing.T) {
		trail := rung.NewRungTrail()
		_, err := trail.UpdateMetaData(tradeId, maxRatingMetadataTrail)
		xrf197ilz35aq0.AssertError(t, err)
		firstUpdatedAt := trail.UpdatedAt()

		updated, err := trail.UpdateMetaData(tradeId, maxRatingMetadataTrail)
		xrf197ilz35aq0.AssertError(t, err)
		assert.False(t, updated)
		assert.Equal(t, trail.UpdatedAt(), firstUpdatedAt)
	})
}

func assertExpectedTrailScore(t testing.TB, r rung.Trail, expectedScore int) {
	t.Helper()
	if r.Score() != expectedScore {
		t.Errorf("r.Score() = %d, want %d", r.Score(), expectedScore)
	}
}

package rung_test

import (
	"testing"
	"xrf197ilz35aq0/model/rung"
)

func TestNewRungTrail(t *testing.T) {
	t.Run("should return a new RungTrail struct", func(t *testing.T) {
		trail := rung.NewRungTrail()
		if trail.Score() != 0 {
			t.Errorf("trail.Score() = %d, want 0", trail.Score())
		}
	})
}

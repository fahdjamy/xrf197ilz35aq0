package rung_test

import (
	"testing"
	"xrf197ilz35aq0/model/rung"
)

func TestRung(t *testing.T) {
	t.Run("creates a new Rung", func(t *testing.T) {
		randomUserFP := "finger-printed"

		newRung := rung.NewRung(randomUserFP)

		if newRung.Magnitude() != 0 {
			t.Errorf("Magnitude should be 0, but was %d", newRung.Magnitude())
		}
	})
}

func TestRungAddTradeTrail(t *testing.T) {
	t.Run("adds a trail to the Rung", func(t *testing.T) {
		randomUserFP := "finger-printed"
		customerFP := "pretty-printed"
		tradeId := "polly"
		newRung := rung.NewRung(randomUserFP)

		newRung.AddTradeTrail(customerFP, tradeId)

		if len(newRung.TradeTrails()) <= 0 {
			t.Errorf("Trails should have been added to the Rung")
		}
	})

	t.Run("adds new customer rung object once", func(t *testing.T) {
		randomUserFP := "finger-printed"
		customerFP := "pretty-printed"
		newRung := rung.NewRung(randomUserFP)
		tradeId := "polly"

		newRung.AddTradeTrail(customerFP, tradeId)
		newRung.AddTradeTrail(customerFP, tradeId)

		if len(newRung.TradeTrails()) > 1 {
			t.Errorf("new TradeRung should have only been added once but was %d", len(newRung.TradeTrails()))
		}
	})
}

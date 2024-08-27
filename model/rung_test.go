package model_test

import (
	"testing"
	"xrf197ilz35aq0/model"
)

func TestRung(t *testing.T) {
	t.Run("creates a new Rung", func(t *testing.T) {
		randomUserFP := "finger-printed"

		rung := model.NewRung(randomUserFP)

		if rung.Magnitude() != 0 {
			t.Errorf("Magnitude should be 0, but was %d", rung.Magnitude())
		}
	})
}

func TestRungAddTradeTrail(t *testing.T) {
	t.Run("adds a trail to the Rung", func(t *testing.T) {
		randomUserFP := "finger-printed"
		customerFP := "pretty-printed"
		tradeId := "polly"
		rung := model.NewRung(randomUserFP)

		rung.AddTradeTrail(customerFP, tradeId)

		if len(rung.TradeTrails()) <= 0 {
			t.Errorf("Trails should have been added to the Rung")
		}
	})

	t.Run("adds new customer rung object once", func(t *testing.T) {
		randomUserFP := "finger-printed"
		customerFP := "pretty-printed"
		rung := model.NewRung(randomUserFP)
		tradeId := "polly"

		rung.AddTradeTrail(customerFP, tradeId)
		rung.AddTradeTrail(customerFP, tradeId)

		if len(rung.TradeTrails()) > 1 {
			t.Errorf("new TradeRung should have only been added once but was %d", len(rung.TradeTrails()))
		}
	})
}

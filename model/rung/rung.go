package rung

import (
	"time"
	"xrf197ilz35aq0/random"
)

// Rung has a 1:1 relationship with the User's table
// there can only be one of this for a user at a time
// it calculates the trust magnitude based off trades with other users
// there should be consent from the buyer (User) to all them to be captured in the rung
type Rung struct {
	id          int64
	magnitude   int
	updated     time.Time
	created     time.Time
	ownerFP     string
	tradeTrails map[string]RungTrail
}

func (r *Rung) AddTradeTrail(customerFP string, tradeId string) {
	tradeTrail, ok := r.tradeTrails[customerFP]
	if ok {
		tradeTrail.UpdateMetaData(tradeId, 0)
	} else {
		rungTrail := NewRungTrail()
		rungTrail.UpdateMetaData(tradeId, 0)
		r.tradeTrails[customerFP] = *rungTrail
	}
	r.updated = time.Now()
}

func (r *Rung) Magnitude() int {
	return r.magnitude
}

func (r *Rung) Updated() time.Time {
	return r.updated
}

func (r *Rung) Created() time.Time {
	return r.created
}

func (r *Rung) TradeTrails() map[string]RungTrail {
	return r.tradeTrails
}

func NewRung(ownerFingerprint string) *Rung {
	now := time.Now()
	id := random.PositiveInt64()
	trails := make(map[string]RungTrail)

	return &Rung{
		id:          id,
		magnitude:   0,
		created:     now,
		updated:     now,
		tradeTrails: trails,
		ownerFP:     ownerFingerprint,
	}
}

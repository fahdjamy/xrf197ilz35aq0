package rung

import (
	"time"
	"xrf197ilz35aq0/internal/random"
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
	tradeTrails map[string]Trail // customer: rungTrails
}

func (r *Rung) AddTradeTrail(customerFP string, tradeId string, metadata TrailMetaData) (bool, error) {
	tradeTrail, ok := r.tradeTrails[customerFP]

	isRungTrailUpdated := false
	if ok {
		updated, err := tradeTrail.UpdateMetaData(tradeId, metadata)
		if err != nil {
			return false, err
		}
		isRungTrailUpdated = updated
	} else {
		rungTrail := NewRungTrail()
		_, err := rungTrail.UpdateMetaData(tradeId, metadata)
		if err != nil {
			return false, err
		}
		r.tradeTrails[customerFP] = *rungTrail
		isRungTrailUpdated = true
	}

	if isRungTrailUpdated {
		r.calculateMagnitude()
		r.updated = time.Now()
	}

	return true, nil
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

func (r *Rung) TradeTrails() map[string]Trail {
	return r.tradeTrails
}

func (r *Rung) calculateMagnitude() {
	if len(r.tradeTrails) > 0 {

	}
}

func NewRung(ownerFingerprint string) *Rung {
	now := time.Now()
	id := random.PositiveInt64()
	trails := make(map[string]Trail)

	return &Rung{
		id:          id,
		magnitude:   0,
		created:     now,
		updated:     now,
		tradeTrails: trails,
		ownerFP:     ownerFingerprint,
	}
}

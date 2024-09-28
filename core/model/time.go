package model

import (
	"encoding/json"
	"time"
	"xrf197ilz35aq0/core"
)

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), core.InternalError{
			Time:    time.Now(),
			Message: "failed to marshal null",
			Source:  "core/model/time#MarshalJSON",
		}
	}
	var err error
	// RFC3339 -> "YYYY-MM-DDTHH:mm:ssZ"
	parsedTime, err := time.Parse(time.RFC3339, t.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	return json.Marshal(parsedTime)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var parsedTime time.Time
	err = json.Unmarshal(data, &parsedTime)
	if err != nil {
		return err
	}
	*t = Time{parsedTime}
	return nil
}

func NewTime(t time.Time) Time {
	return Time{t}
}

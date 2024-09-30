package model

import (
	"encoding/json"
	"time"
	error2 "xrf197ilz35aq0/internal/error"
)

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return t.Time.Format(time.RFC3339)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), &error2.Internal{
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

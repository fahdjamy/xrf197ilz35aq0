package xrf197ilz35aq0

import "encoding/json"

type SerializableString string

func (s *SerializableString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*s))
}

func (s *SerializableString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = SerializableString(str)
	return nil
}

package xrf197ilz35aq0

type Serializable interface {
	// MarshalJSON Takes a Go data structure (like a struct, map, or slice)
	// and converts it into a JSON-formatted string.
	// Used when sending data over a network, store it in a file,
	//	or communicate with external systems that expect JSON.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON Takes a JSON-formatted string and converts it back into a Go data structure.
	// This allows to work with data that you've received in JSON format and manipulate it using familiar Go constructs.
	UnmarshalJSON([]byte) error
}

package encryption

type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

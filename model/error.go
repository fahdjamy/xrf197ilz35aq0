package model

type InvalidRequest struct {
	Message string
}

func (i InvalidRequest) Error() string {
	return i.Message
}

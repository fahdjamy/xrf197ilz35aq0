package core

import (
	"fmt"
	"time"
)

type InvalidRequest struct {
	Message string
	Err     error
	Time    time.Time
}

func (i InvalidRequest) Error() string {
	return i.Message
}

type InternalError struct {
	Message string
	Time    time.Time
	Source  string
	Err     error
}

func (i InternalError) Error() string {
	return fmt.Sprintf("message=%s :: time=%s :: source%s :: err=%s", i.Message, i.Time, i.Source, i.Err)
}

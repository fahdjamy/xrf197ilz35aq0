package error

import (
	"fmt"
	"time"
)

type External struct {
	Err     error
	Source  string
	Message string
	Time    time.Time
}

func (e *External) WithErr(msg string, err error) *External {
	e.Err = err
	e.Message = msg
	return e
}

func (e *External) NoErr(msg string) *External {
	e.Message = msg
	return e
}

func (e *External) Error() string {
	return fmt.Sprintf("message=%s :: time=%s :: source%s :: \n\terr=%s", e.Message, e.Time, e.Source, e.Err)
}

type Internal struct {
	Message string
	Time    time.Time
	Source  string
	Err     error
}

func (ie *Internal) NoErr(msg string) *Internal {
	ie.Message = msg
	return ie
}

func (ie *Internal) WithErr(msg string, err error) *Internal {
	ie.Err = err
	ie.Message = msg
	return ie
}

func (ie *Internal) Error() string {
	return fmt.Sprintf("message=%s :: time=%s :: source%s :: \n\terr=%s", ie.Message, ie.Time, ie.Source, ie.Err)
}

package error

import (
	"fmt"
)

type External struct {
	Err     error
	Source  string
	Message string
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
	return fmt.Sprintf("%s", e.Message)
}

func (e *External) String() string {
	str := fmt.Sprintf("message=%s :: source%s", e.Message, e.Source)
	if e.Err != nil {
		str += fmt.Sprintf(" :: \n\t%s", e.Err)
	}
	return str
}

type Internal struct {
	Message string
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
	return fmt.Sprintf("%s", ie.Message)
}

func (e *Internal) String() string {
	str := fmt.Sprintf("message=%s :: source%s", e.Message, e.Source)
	if e.Err != nil {
		str += fmt.Sprintf(" :: \n\t%s", e.Err)
	}
	return str
}

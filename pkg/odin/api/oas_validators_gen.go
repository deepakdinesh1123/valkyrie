// Code generated by ogen, DO NOT EDIT.

package api

import (
	"github.com/go-faster/errors"
)

func (s GetExecutionResultsOKApplicationJSON) Validate() error {
	alias := ([]ExecutionResult)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s GetExecutionsOKApplicationJSON) Validate() error {
	alias := ([]Execution)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

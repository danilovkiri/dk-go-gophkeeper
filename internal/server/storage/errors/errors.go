// Package errors provides custom error types.
package errors

import (
	"fmt"
)

type (
	WrongDBError struct {
		Err error
		ID  string
	}
	StatementPSQLError struct {
		Err error
	}
	AlreadyExistsError struct {
		Err error
		ID  string
	}
	ExecutionPSQLError struct {
		Err error
	}
	ContextTimeoutExceededError struct {
		Err error
	}
	NotFoundError struct {
		Err error
	}
	ScanningPSQLError struct {
		Err error
	}
	InvalidPasswordError struct {
		Err error
	}
)

func (e *WrongDBError) Error() string {
	return fmt.Sprintf("%s: invalid DB indetifier", e.ID)
}

func (e *StatementPSQLError) Error() string {
	return fmt.Sprintf("%s: could not compile", e.Err.Error())
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s: already exists", e.ID)
}

func (e *ExecutionPSQLError) Error() string {
	return fmt.Sprintf("%s: could not execute", e.Err.Error())
}

func (e *ContextTimeoutExceededError) Error() string {
	return fmt.Sprintf("%s: context timeout exceeded", e.Err.Error())
}

func (e *NotFoundError) Error() string {
	return "not found in storage"
}

func (e *ScanningPSQLError) Error() string {
	return fmt.Sprintf("%s: could not scan rows", e.Err.Error())
}

func (e *InvalidPasswordError) Error() string {
	return "password is invalid"
}

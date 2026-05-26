package errorspkg

import (
	"fmt"
	"reflect"
	"strings"
)

type NotFoundError struct {
	what string
	by   []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found by [%s]", e.what, strings.Join(e.by, ", "))
}

func NewNotFoundError(what string, by ...string) error {
	return NotFoundError{
		what: what,
		by:   by,
	}
}

type NotValidError struct {
	what string
	need string
}

func (e NotValidError) Error() string {
	return fmt.Sprintf("%s is not valid %s", e.what, e.need)
}

func NewNotValidError(what, need string) error {
	return NotValidError{
		what: what,
		need: need,
	}
}

type IsEmptyError struct {
	what string
}

func (e IsEmptyError) Error() string {
	return fmt.Sprintf("%s is empty", e.what)
}

func NewIsEmptyError(what string) error {
	return IsEmptyError{
		what: what,
	}
}

type FailedError struct {
	what string
}

func (e FailedError) Error() string {
	return fmt.Sprintf("filed to %s", e.what)
}

func NewFailedError(what string) error {
	return FailedError{
		what: what,
	}
}

type ReadBodyError struct {
	err error
}

func (e ReadBodyError) Error() string {
	return fmt.Sprintf("failed to read body: %s", e.err.Error())
}

func NewReadBodyError(err error) error {
	return ReadBodyError{
		err: err,
	}
}

type ValidationError struct {
	Constructor string
	StructName  string
	Err         error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for constructor [%s] struct [%s]: %v",
		e.Constructor,
		e.StructName,
		e.Err,
	)
}

func NewValidationError(constructor string, obj any, err error) error {
	var structName string

	if obj != nil {
		structName = reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
	}

	return &ValidationError{
		Constructor: constructor,
		StructName:  structName,
		Err:         err,
	}
}

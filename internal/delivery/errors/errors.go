package errors

import "github.com/pkg/errors"

var (
	ErrBadInput     = errors.New("bad input")
	ErrWrongData    = errors.New("wrong data")
	ErrAccessDenied = errors.New("access denied")
	ErrNotAllowed   = errors.New("not allowed")
	ErrInternal     = errors.New("internal error")
	ErrNotFound     = errors.New("not found")
)

func NewErrBadInput(msg string) error {
	return errors.Wrap(ErrBadInput, msg)
}

func NewErrWrongData(msg string) error {
	return errors.Wrap(ErrWrongData, msg)
}

func NewErrAccessDenied(msg string) error {
	return errors.Wrap(ErrAccessDenied, msg)
}

func NewErrInternal(msg string) error {
	return errors.Wrap(ErrInternal, msg)
}

func NewErrNotFound(msg string) error {
	return errors.Wrap(ErrNotFound, msg)
}

func NewErrNotAllowed(msg string) error {
	return errors.Wrap(ErrNotAllowed, msg)
}

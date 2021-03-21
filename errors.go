package errors

import (
	stderrors "errors"
)

// New provides an interface to the package taking a single generic error.
func New(err error) error {
	return Wrap(err, nil)
}

// Wrap provides an interface to the package that takes two generic errors,
// wrapping one inside the other.
func Wrap(err, cause error) error {
	if err == nil {
		return nil
	}

	return &wrappedError{
		cause: cause,
		err:   err,
	}
}

// wrappedError represents an error type that can be used to wrap errors.
type wrappedError struct {
	cause error
	err   error
}

func (e *wrappedError) Error() string {
	if e.cause == nil {
		return e.err.Error()
	}

	return e.err.Error() + ": " + e.cause.Error()
}

// Unwrap returns the wrapped error, or nil if the wrapped error is unset.
func (e *wrappedError) Unwrap() error {
	return e.cause
}

// Is reports whether any error in e's chain matches target.
// https://github.com/golang/go/blob/7bb721b9384bdd196befeaed593b185f7f2a5589/src/errors/wrap.go#L24-L39
func (e *wrappedError) Is(target error) bool {
	if stderrors.Is(e.err, target) {
		return true
	}

	return stderrors.Is(e.cause, target)
}

// As finds the first error in e's chain that matches target, and if so, sets
// target to that error value and returns true.
// https://github.com/golang/go/blob/7bb721b9384bdd196befeaed593b185f7f2a5589/src/errors/wrap.go#L61-L77
func (e *wrappedError) As(target interface{}) bool {
	if stderrors.As(e.err, target) {
		return true
	}

	return stderrors.As(e.cause, target)
}

// NewWithMessage provides an interface to the package taking a single generic
// error and a contextual message.
func NewWithMessage(err error, msg string) error {
	return WrapWithMessage(err, nil, msg)
}

// WrapWithMessage provides an interface to the package that takes two generic
// errors, wrapping one inside the other, and a contextual message.
func WrapWithMessage(err, cause error, msg string) error {
	if err == nil {
		return nil
	}

	return &withMessage{
		wrappedError: &wrappedError{
			cause: cause,
			err:   err,
		},
		msg: msg,
	}
}

type withMessage struct {
	*wrappedError
	msg string
}

func (e *withMessage) Error() string {
	if e.cause == nil {
		return e.msg + ": " + e.err.Error()
	}

	return e.msg + ": " + e.err.Error() + ": " + e.cause.Error()
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type
// contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise it returns false.
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}
